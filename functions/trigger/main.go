package trigger

import (
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
	shooterURLs := []string{
		"https://us-central1-a-global-presence-hackattic-db.cloudfunctions.net/shooter",
		"https://us-central1-a-global-presence-hackattic-db.cloudfunctions.net/shooter",
		"https://us-central1-a-global-presence-hackattic-db.cloudfunctions.net/shooter",
	}

	docs := client.Collection(cfg.ChallengeDoc).Where("status", "!=",
		"completed").OrderBy("status", firestore.Desc).Documents(ctx)

	doc, err := docs.Next()

	var firestoreDoc hackattic.Document

	err = doc.DataTo(&attempt)

	if err != nil {
		cfg.Logger.Error(fmt.Sprintf("Failed to decode document: %v\n", err))
		sendErr(w, "Failed to decode document", []error{err})
	}

	go func() {
		for _, url := range shooterURLs {
			shooterRespo, err := http.Get(url)
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
			})
		}

		return
	}()

	// Check responses

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResp := map[string]string{
		"status":     "Document created successfully",
		"documentId": result.ID,
	}

	json.NewEncoder(w).Encode(jsonResp)
	return
}
