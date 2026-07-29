package handlers

import (
	"encoding/json"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthReturnsStatusOK(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rr := httptest.NewRecorder()

	Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rr.Code,
		)
	}
	contentType := rr.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Fatalf(
			"expected application/json, got %s",
			contentType,
		)
	}

	var response HealthResponse

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Status != "ok" {
		t.Fatalf(
			"expected status ok, got %s",
			response.Status,
		)
	}
	if response.Agent.Name != version.Build.Name {
		t.Fatalf(
			"expected agent name %q, got %q",
			version.Build.Name,
			response.Agent.Name,
		)
	}

	if response.Agent.Version != version.Build.Version {
		t.Fatalf(
			"expected version %q, got %q",
			version.Build.Version,
			response.Agent.Version,
		)
	}

}
