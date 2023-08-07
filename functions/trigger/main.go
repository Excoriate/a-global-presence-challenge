package trigger

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"github.com/Excoriate/a-global-presence-challenge/pkg/hackattic"
	"github.com/Excoriate/a-global-presence-challenge/pkg/utils"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"net/http"
	"os"
)

func init() {
	functions.HTTP("trigger", trigger)
}

func getPresences() []string {
	return []string{"europe-west4", "us-central1", "asia-east1", "australia-southeast1", "europe-north1", "europe-west1", "northamerica-northeast1", "asia-northeast1", "asia-southeast1"}
}

func sendErr(w http.ResponseWriter, status string, errs []error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	var errMsgs []string
	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	jsonResp := map[string]interface{}{
		"status": status,
		"errors": errMsgs,
	}

	json.NewEncoder(w).Encode(jsonResp)
}

func trigger(w http.ResponseWriter, r *http.Request) {
	// Getting and building the configuration.
	cfg, errs := config.New().
		WithEnv(os.Getenv("ENVIRONMENT")).
		WithScannedEnvVars().
		WithRequiredConfig().
		Build()

	if errs != nil && len(errs) > 0 {
		sendErr(w, "Failed to build config", errs)
		return
	}

	// Configure the Challenge attempt.
	challengeClient := hackattic.New(cfg)
	challengeCfg, errs := challengeClient.
		WithHeaders([]string{}).
		WithHTTPClient().
		WithAccessToken().
		Build()

	if errs != nil && len(errs) > 0 {
		sendErr(w, "Failed to build challenge configuration", errs)
		return
	}

	// Configure a new attempt.
	attempt := hackattic.NewAttempt(cfg)
	attemptCfg, errs := attempt.
		WithChallenge(challengeCfg).
		WithNewPresenceToken().
		Build()

	if errs != nil && len(errs) > 0 {
		sendErr(w, "Failed to build attempt configuration", errs)
		return
	}

	cfg.Logger.Info(fmt.Sprintf("New attempt with presence token: %v\n", attemptCfg.PresenceToken))
	ctx := context.Background()

	// Firestore client
	client, dbErr := firestore.NewClient(ctx, cfg.ProjectId)
	if dbErr != nil {
		sendErr(w, "Failed to create Firestore client", []error{dbErr})
		return
	}

	defer client.Close()

	result, _, docErr := client.Collection(cfg.ChallengeDoc).Add(ctx, map[string]interface{}{
		"attempt_id":    utils.GetUUID(),
		"accessToken":   challengeCfg.AccessToken,
		"presenceToken": attemptCfg.PresenceToken,
		"attempts": []map[string]interface{}{
			{
				"attempt_id":            utils.GetUUID(),
				"presenceToken":         attemptCfg.PresenceToken,
				"isCompleted":           false,
				"status":                "trigger-init",
				"hackatticCountryCheck": "",
				"shooter_url":           "",
			},
		},
		"status":           "pending",
		"numberOfAttempts": 0,
		"result":           "",
	})

	if docErr != nil {
		sendErr(w, "Failed to create document", []error{docErr})
		return
	}

	// Call each shooter in a goroutine.
	shooterURLs := []string{}
	for _, region := range getPresences() {
		shooterURLs = append(shooterURLs, fmt.Sprintf("https://%v-%v.cloudfunctions."+
			"net/a-global-presence-function-shooter", region, cfg.ProjectId))
	}

	docs := client.Collection(cfg.ChallengeDoc).Where("status", "!=",
		"completed").OrderBy("status", firestore.Desc).Documents(ctx)

	doc, err := docs.Next()

	var firestoreDoc hackattic.Document
	err = doc.DataTo(&firestoreDoc)

	if err != nil {
		cfg.Logger.Error(fmt.Sprintf("Failed to decode document: %v\n", err))
		sendErr(w, "Failed to decode document", []error{err})
	}

	counterSuccess := 0
	for _, url := range shooterURLs {
		// Expected response
		//{"HackatticCountryCheck":"US","IsSuccessful":true,"Status":"checked","Error":null}
		//{"response":{"HackatticCountryCheck":"US","IsSuccessful":true,"Status":"checked","Error":null},"status":"success"}
		shooterRespo, err := challengeCfg.HttpClient.Get(url)
		if err != nil {
			cfg.Logger.Error(fmt.Sprintf("Failed to call shooter: %v\n", err))
		}

		defer shooterRespo.Body.Close()
		cfg.Logger.Info(fmt.Sprintf("Shooter response: %v\n", shooterRespo))

		responseJson := hackattic.ShooterResponse{}
		err = json.NewDecoder(shooterRespo.Body).Decode(&responseJson)
		if err != nil {
			cfg.Logger.Error(fmt.Sprintf("Failed to decode shooter response: %v\n", err))
		}

		// Update the document
		firestoreDoc.Attempts = append(firestoreDoc.Attempts, hackattic.AttemptRegisterDoc{
			AttemptID:             utils.GetUUID(),
			HackatticCountryCheck: responseJson.HackatticCountryCheck,
			IsCompleted:           true,
			PresenceToken:         attemptCfg.PresenceToken,
			Status:                "hackattic-checked-ok",
			ShooterURl:            url,
		})

		// get new presence token
		presenceTokenNew, _ := challengeCfg.HttpClient.Get(challengeCfg.APIGetPresenceToken)
		defer presenceTokenNew.Body.Close()
		presenteTOkenNewJson := hackattic.PresenceTokenResponse{}
		err = json.NewDecoder(presenceTokenNew.Body).Decode(&presenteTOkenNewJson)
		if err != nil {
			cfg.Logger.Error(fmt.Sprintf("Failed to decode presence token: %v\n", err))
		}

		_, err = client.Collection(cfg.ChallengeDoc).Doc(result.ID).Set(ctx, firestoreDoc)
		if err != nil {
			cfg.Logger.Error(fmt.Sprintf("Failed to update document: %v\n", err))
			sendErr(w, "Failed to update document", []error{err})
		}

		if responseJson.IsSuccessful {
			counterSuccess++
		}
	}

	// Make the final attempt to send an empty json {}
	if counterSuccess >= 7 {
		resp, err := challengeCfg.HttpClient.Post(challengeCfg.APISubmitSolution, "application/json", bytes.NewBuffer([]byte("{}")))
		if err != nil {
			cfg.Logger.Error(fmt.Sprintf("Failed to submit solution: %v\n", err))
			sendErr(w, "Failed to submit solution", []error{err})
		}

		defer resp.Body.Close()
		cfg.Logger.Info(fmt.Sprintf("Submit solution response: %v\n", resp))

		// Check responses
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jsonResp := map[string]string{
			"status":     "The challenge is completed",
			"documentId": result.ID,
		}

		json.NewEncoder(w).Encode(jsonResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResp := map[string]string{
		"status":     "The challenge is not completed",
		"documentId": result.ID,
	}

	json.NewEncoder(w).Encode(jsonResp)
	return
}
