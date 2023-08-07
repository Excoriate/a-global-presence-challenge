package shooter

import (
	"encoding/json"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"github.com/Excoriate/a-global-presence-challenge/pkg/hackattic"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"net/http"
	"os"
)

func init() {
	functions.HTTP("shooter", shooter)
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

func shooter(w http.ResponseWriter, r *http.Request) {
	presenceTokenReceived := r.URL.Query().Get("presence_token")
	if presenceTokenReceived == "" {
		sendErr(w, "No presence token received", []error{})
		return
	}

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

	cfg.Logger.Info("Presence token received: " + presenceTokenReceived)

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
		//WithNewPresenceToken(). isn't required at the Shooter level. It's already in the DB.
		WithChallenge(challengeCfg).
		Build()

	if errs != nil && len(errs) > 0 {
		sendErr(w, "Failed to build attempt configuration", errs)
		return
	}

	// New Shooter
	shooterInst := hackattic.NewShooter()
	shooterCfg, errs := shooterInst.
		WithConfig(cfg).
		WithChallenge(challengeCfg).
		WithAttempt(attemptCfg).
		WithPassedPresenceToken(presenceTokenReceived). // Here we're obtaining the presence token from the DB.
		WithCountryCheck(). // Calling the country check 'endpoint'
		Complete()

	if errs != nil && len(errs) > 0 {
		sendErr(w, "Failed to build shooter configuration", errs)
		return
	}

	//result := shooterCfg.Register()
	response := shooterCfg.Response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return response
	jsonResp := map[string]interface{}{
		"status":   "success",
		"response": response,
	}

	json.NewEncoder(w).Encode(jsonResp)
	return
}
