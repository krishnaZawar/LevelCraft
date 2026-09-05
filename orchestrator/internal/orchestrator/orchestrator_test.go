package orchestrator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

func Test_NewOrchestrator(t *testing.T) {
	o := NewOrchestrator()

	if o == nil {
		t.Fatal("NewOrchestrator() returned nil")
	}

	if o.processes == nil {
		t.Fatal("processes map should be initialized")
	}

	if len(o.processes) != 0 {
		t.Fatalf("expected empty processes map, got %d processes", len(o.processes))
	}
}

func Test_RunProcess(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		o := NewOrchestrator()

		expectedProcess := &entity.Process{
			Name:             base.Process_EditorBackend,
			CommunicationURI: "http://localhost:8080",
		}

		fn := func() (*entity.Process, error) {
			return expectedProcess, nil
		}

		err := o.runProcess(fn)
		if err != nil {
			t.Fatalf("runProcess() returned error: %v", err)
		}

		process, ok := o.processes[base.Process_EditorBackend]
		if !ok {
			t.Fatal("process was not added to orchestrator")
		}

		if process != expectedProcess {
			t.Fatalf("stored process = %p, want %p", process, expectedProcess)
		}
	})

	t.Run("unsuccessful creation", func(t *testing.T) {
		o := NewOrchestrator()

		expectedErr := context.Canceled

		fn := func() (*entity.Process, error) {
			return nil, expectedErr
		}

		err := o.runProcess(fn)

		if err != expectedErr {
			t.Fatalf("runProcess() error = %v, want %v", err, expectedErr)
		}

		if len(o.processes) != 0 {
			t.Fatalf("expected no processes to be registered, got %d", len(o.processes))
		}
	})
}

func Test_Monitor(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		closeServer bool
		want        applicationState
	}{
		{
			name:       "healthy editor backend",
			statusCode: http.StatusOK,
			want:       applicationState_healthy,
		},
		{
			name:       "editor backend unhealthy",
			statusCode: http.StatusInternalServerError,
			want:       applicationState_editorExit,
		},
		{
			name:       "editor backend not found",
			statusCode: http.StatusNotFound,
			want:       applicationState_editorExit,
		},
		{
			name:        "editor backend unavailable",
			closeServer: true,
			want:        applicationState_editorExit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))

			if tt.closeServer {
				server.Close()
			} else {
				defer server.Close()
			}

			o := NewOrchestrator()

			o.processes[base.Process_EditorBackend] = &entity.Process{
				Name:             base.Process_EditorBackend,
				CommunicationURI: server.URL,
			}

			got := o.monitor()

			if got != tt.want {
				t.Errorf("monitor() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_ExitProcesses(t *testing.T) {
	t.Run("no missing processes", func(t *testing.T) {
		cmd := exec.Command("sleep", "10")

		if err := cmd.Start(); err != nil {
			t.Fatalf("failed to start test process: %v", err)
		}

		o := NewOrchestrator()

		process := &entity.Process{
			Name:             base.Process_EditorBackend,
			Cmd:              cmd,
			CommunicationURI: "http://localhost:8080",
		}

		o.processes[base.Process_EditorBackend] = process

		o.exitProcesses([]string{
			base.Process_EditorBackend,
		})

		if _, ok := o.processes[base.Process_EditorBackend]; ok {
			t.Fatal("process should have been removed from processes map")
		}

		// Make sure the process has actually exited.
		done := make(chan error, 1)
		go func() {
			_, err := cmd.Process.Wait()
			done <- err
		}()

		select {
		case <-done:
			// Process exited successfully.
		case <-time.After(time.Second):
			t.Fatal("process did not exit after Stop()")
		}
	})

	t.Run("missing processes", func(t *testing.T) {
		o := NewOrchestrator()

		// Should not panic when the process doesn't exist.
		o.exitProcesses([]string{
			base.Process_EditorBackend,
		})

		if len(o.processes) != 0 {
			t.Fatalf("expected empty processes map, got %d", len(o.processes))
		}
	})
}
