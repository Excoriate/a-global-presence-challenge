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

func sendErr(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	jsonResp := map[string]string{
		"status": errMsg,
	}

	json.NewEncoder(w).Encode(jsonResp)
}

func trigger(w http.ResponseWriter, r *http.Request) {
	// Getting and building the configuration.
	cfg, err := config.New().
		WithEnv(os.Getenv("ENVIRONMENT")).
		WithScannedEnvVars().
		WithRequiredConfig().
		Build()

	if err != nil {
		sendErr(w, fmt.Sprintf("Failed to build config: %v\n", err))
		return
	}

	// Configure the Challenge attempt.
	challengeClient := hackattic.New(cfg)
	challengeCfg, err := challengeClient.
		WithHeaders([]string{}).
		WithHTTPClient().
		WithAccessToken().
		Build()

	if err != nil {
		sendErr(w, fmt.Sprintf("Failed to build challenge config: %v\n", err))
		return
	}

	// Configure a new attempt.
	attempt := hackattic.NewAttempt(cfg)
	attemptCfg, err := attempt.
		WithChallenge(challengeCfg).
		WithNewPresenceToken().
		Build()

	if err != nil {
		sendErr(w, fmt.Sprintf("Failed to build attempt config: %v\n", err))
		return
	}

	cfg.Logger.Info(fmt.Sprintf("New attempt with presence token: %v\n", attemptCfg.PresenceToken))
	ctx := context.Background()

	// Firestore client
	client, err := firestore.NewClient(ctx, cfg.ProjectId)
	if err != nil {
		sendErr(w, fmt.Sprintf("Failed to create firestore client: %v\n", err))
		return
	}

	defer client.Close()

	result, _, err := client.Collection(cfg.ChallengeDoc).Add(ctx, map[string]interface{}{
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

	if err != nil {
		sendErr(w, fmt.Sprintf("Failed to create document: %v\n", err))
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
