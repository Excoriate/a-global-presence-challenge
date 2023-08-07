package shooter

import (
	"encoding/json"
	"fmt"
	"github.com/Excoriate/a-global-presence-challenge/pkg/config"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"net/http"
	"os"
)

func init() {
	functions.HTTP("shooter", shooter)
}

func shooter(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.New().WithEnv(os.Getenv("ENV")).WithRequiredConfig().Build()

	if err != nil {
		errMsg := fmt.Sprintf("Failed to build config: %v\n", err)
		cfg.Logger.Error(errMsg)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		jsonResp := map[string]string{
			"status": errMsg,
		}

		json.NewEncoder(w).Encode(jsonResp)
		return
	}

	// Get the token from the request as it was passed as a query string param
	token := r.URL.Query().Get("token")
	if token == "" {
		errMsg := fmt.Sprintf("No token provided")
		cfg.Logger.Error(errMsg)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		jsonResp := map[string]string{
			"status": errMsg,
		}

		json.NewEncoder(w).Encode(jsonResp)
		return
	}
}
