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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResp := map[string]string{
		"status":     "Document created successfully",
		"documentId": result.ID,
	}

	json.NewEncoder(w).Encode(jsonResp)
	return
}
