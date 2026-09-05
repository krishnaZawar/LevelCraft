package executor

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_GetRandomPort(t *testing.T) {
	for i := 0; i < 100; i++ {
		portStr := getRandomPort()

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatalf("getRandomPort() returned invalid port %q: %v", portStr, err)
		}

		if port < base.MinPortValue || port > base.MaxPortValue {
			t.Fatalf(
				"getRandomPort() = %d, want value between %d and %d",
				port,
				base.MinPortValue,
				base.MaxPortValue,
			)
		}
	}
}

func Test_CheckProcessHealth(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "healthy",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "unhealthy",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/ping" {
					t.Errorf("expected /ping, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			err := CheckProcessHealth(server.URL)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckProcessHealth() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func Test_CheckProcessHealth_ServerUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()

	err := CheckProcessHealth(url)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func Test_BuildAndRunProcess(t *testing.T) {
	t.Run("successful creation of process", func(t *testing.T) {
		name := "lsCmd"
		cmd := &entity.CommandConfig{
			Pwd:  ".",
			Name: "ls",
			Args: []string{},
			Port: "1010",
		}
		process, err := buildAndRunProcess(name, cmd)
		assert.Nil(t, err)
		assert.Equal(t, name, process.Name)
		assert.Equal(t, "http://localhost:1010", process.CommunicationURI)

		_ = process.Stop()
	})

	t.Run("unsuccessful creation of process", func(t *testing.T) {
		name := "lsCmd"
		cmd := &entity.CommandConfig{
			Pwd:  ".",
			Name: "l",
			Args: []string{},
			Port: "1010",
		}
		process, err := buildAndRunProcess(name, cmd)
		assert.NotNil(t, err)
		assert.Nil(t, process)
	})
}
