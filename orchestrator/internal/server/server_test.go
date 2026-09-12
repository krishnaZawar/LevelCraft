package server

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

// The editor is handed this address at startup, so it has to be serving by
// the time Start returns.
func Test_Start_ReturnsAReachableAddress(t *testing.T) {
	baseURL, err := Start(&fakeController{})
	if err != nil {
		t.Fatalf("Start() returned error: %v", err)
	}

	if !strings.HasPrefix(baseURL, "http://localhost:") {
		t.Fatalf("Start() = %q, want a http://localhost:<port> address", baseURL)
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/ping")
	if err != nil {
		t.Fatalf("failed to reach the control API at %s: %v", baseURL, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ping status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// A port collision would leave an editor talking to the wrong orchestrator.
func Test_Start_AssignsDistinctPorts(t *testing.T) {
	first, err := Start(&fakeController{})
	if err != nil {
		t.Fatalf("Start() returned error: %v", err)
	}
	second, err := Start(&fakeController{})
	if err != nil {
		t.Fatalf("Start() returned error: %v", err)
	}

	if first == second {
		t.Fatalf("both servers bound the same address: %s", first)
	}
}
