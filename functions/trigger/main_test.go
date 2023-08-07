package trigger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTrigger(t *testing.T) {
	req, err := http.NewRequest("GET", "/trigger", nil)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("ENVIRONMENT", "sandbox")
	os.Setenv("CHALLENGE_DB_NAME", "a-global-presence-hackattic-db")
	os.Setenv("CHALLENGE_DOC_NAME", "challenge_doc")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(trigger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
