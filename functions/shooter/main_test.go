package shooter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestShooter(t *testing.T) {
	// Create an initial request to setup.
	req, err := http.NewRequest("GET", "/shooter?presence_token=TEST", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Calling the handler with a query param presence_token", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "sandbox")
		os.Setenv("CHALLENGE_DB_NAME", "a-global-presence-hackattic-db")
		os.Setenv("CHALLENGE_DOC_NAME", "challenge_doc")

		rr := httptest.NewRecorder()

		// Add presence_token as a query parameter.
		q := req.URL.Query()
		q.Add("presence_token", "test")
		req.URL.RawQuery = q.Encode()

		handler := http.HandlerFunc(shooter)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
