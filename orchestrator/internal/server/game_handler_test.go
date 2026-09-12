package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

// records what the handlers asked for, so the control API can be tested
// without spawning any real process
type fakeController struct {
	runCalled  bool
	runPath    string
	runErr     error
	stopCalled bool
	running    bool
}

func (f *fakeController) RunGame(scenePath string) error {
	f.runCalled = true
	f.runPath = scenePath
	return f.runErr
}

func (f *fakeController) StopGame() {
	f.stopCalled = true
}

func (f *fakeController) IsGameRunning() bool {
	return f.running
}

func newTestApp(controller GameController) *fiber.App {
	app := fiber.New()
	CreateRoutes(app, controller)
	return app
}

func doRequest(t *testing.T, app *fiber.App, method string, path string, body string) (*http.Response, []byte) {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = bytes.NewBufferString(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	return resp, data
}

func Test_HandleRunGame(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		runErr     error
		wantStatus int
		wantRun    bool
	}{
		{
			name:       "valid request",
			body:       `{"filepath":"/tmp/scene.json"}`,
			wantStatus: http.StatusOK,
			wantRun:    true,
		},
		{
			name:       "missing filepath",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "non json scene file",
			body:       `{"filepath":"/tmp/scene.txt"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed body",
			body:       `not json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "orchestrator failed to start the game",
			body:       `{"filepath":"/tmp/scene.json"}`,
			runErr:     errors.New("builder backend did not become healthy"),
			wantStatus: http.StatusInternalServerError,
			wantRun:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &fakeController{runErr: tt.runErr}
			resp, body := doRequest(t, newTestApp(controller), http.MethodPost, "/game/run", tt.body)

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", resp.StatusCode, tt.wantStatus, body)
			}

			if controller.runCalled != tt.wantRun {
				t.Fatalf("RunGame called = %v, want %v", controller.runCalled, tt.wantRun)
			}

			var parsed entity.GameActionResponse
			if err := json.Unmarshal(body, &parsed); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}
			if parsed.Success != (tt.wantStatus == http.StatusOK) {
				t.Fatalf("success = %v, want %v", parsed.Success, tt.wantStatus == http.StatusOK)
			}
		})
	}
}

// No process should be spawned for a scene file the builder could never load.
func Test_HandleRunGame_DoesNotSpawnForInvalidPaths(t *testing.T) {
	controller := &fakeController{}
	resp, _ := doRequest(t, newTestApp(controller), http.MethodPost, "/game/run", `{"filepath":"scene.yaml"}`)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
	if controller.runCalled {
		t.Fatal("RunGame was called for a non json scene file")
	}
}

func Test_HandleStopGame(t *testing.T) {
	controller := &fakeController{}
	resp, _ := doRequest(t, newTestApp(controller), http.MethodPost, "/game/stop", "")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if !controller.stopCalled {
		t.Fatal("StopGame was not called")
	}
}

func Test_HandleGameStatus(t *testing.T) {
	for _, running := range []bool{true, false} {
		controller := &fakeController{running: running}
		resp, body := doRequest(t, newTestApp(controller), http.MethodGet, "/game/status", "")

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var parsed entity.GameStatusResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}
		if parsed.Running != running {
			t.Fatalf("running = %v, want %v", parsed.Running, running)
		}
	}
}

func Test_Ping(t *testing.T) {
	resp, body := doRequest(t, newTestApp(&fakeController{}), http.MethodGet, "/ping", "")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if !strings.Contains(string(body), "pong") {
		t.Fatalf("unexpected ping response: %s", body)
	}
}
