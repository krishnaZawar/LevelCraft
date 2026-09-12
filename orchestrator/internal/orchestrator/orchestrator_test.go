package orchestrator

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"sync"
	"sync/atomic"
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

func Test_Monitor_BuilderProcesses(t *testing.T) {
	tests := []struct {
		name        string
		processName string
		statusCode  int
		closeServer bool
		want        applicationState
	}{
		{
			name:        "healthy builder backend",
			processName: base.Process_BuilderBackend,
			statusCode:  http.StatusOK,
			want:        applicationState_healthy,
		},
		{
			name:        "builder backend unhealthy",
			processName: base.Process_BuilderBackend,
			statusCode:  http.StatusInternalServerError,
			want:        applicationState_builderExit,
		},
		{
			name:        "builder frontend unavailable",
			processName: base.Process_BuilderFrontend,
			closeServer: true,
			want:        applicationState_builderExit,
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
			o.processes[tt.processName] = &entity.Process{
				Name:             tt.processName,
				CommunicationURI: server.URL,
			}

			if got := o.monitor(); got != tt.want {
				t.Errorf("monitor() = %q, want %q", got, tt.want)
			}
		})
	}
}

// An editor going down takes the application with it, so it is the failure
// that has to be reported when both happen in the same pass.
func Test_Monitor_EditorFailureOutranksBuilderFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	o := NewOrchestrator()
	o.processes[base.Process_BuilderBackend] = &entity.Process{
		Name:             base.Process_BuilderBackend,
		CommunicationURI: server.URL,
	}
	o.processes[base.Process_EditorBackend] = &entity.Process{
		Name:             base.Process_EditorBackend,
		CommunicationURI: server.URL,
	}

	if got := o.monitor(); got != applicationState_editorExit {
		t.Errorf("monitor() = %q, want %q", got, applicationState_editorExit)
	}
}

func Test_IsGameRunning(t *testing.T) {
	t.Run("no builder processes", func(t *testing.T) {
		o := NewOrchestrator()
		o.processes[base.Process_EditorBackend] = &entity.Process{Name: base.Process_EditorBackend}

		if o.IsGameRunning() {
			t.Fatal("IsGameRunning() = true with only editor processes tracked, want false")
		}
	})

	t.Run("builder process tracked", func(t *testing.T) {
		o := NewOrchestrator()
		o.processes[base.Process_BuilderBackend] = &entity.Process{Name: base.Process_BuilderBackend}

		if !o.IsGameRunning() {
			t.Fatal("IsGameRunning() = false with a builder process tracked, want true")
		}
	})
}

func Test_RunGame_RejectsASecondRun(t *testing.T) {
	o := NewOrchestrator()
	o.processes[base.Process_BuilderBackend] = &entity.Process{Name: base.Process_BuilderBackend}

	err := o.RunGame("scene.json")

	if !errors.Is(err, ErrGameAlreadyRunning) {
		t.Fatalf("RunGame() error = %v, want %v", err, ErrGameAlreadyRunning)
	}

	// The already-running game must survive a rejected second request.
	if !o.IsGameRunning() {
		t.Fatal("the running game was torn down by a rejected second run")
	}
}

func Test_StopGame(t *testing.T) {
	t.Run("stops builder processes and leaves the editor alone", func(t *testing.T) {
		cmd := exec.Command("sleep", "10")
		if err := cmd.Start(); err != nil {
			t.Fatalf("failed to start test process: %v", err)
		}

		o := NewOrchestrator()
		o.processes[base.Process_BuilderBackend] = &entity.Process{
			Name: base.Process_BuilderBackend,
			Cmd:  cmd,
		}
		o.processes[base.Process_EditorBackend] = &entity.Process{
			Name: base.Process_EditorBackend,
		}

		o.StopGame()

		if _, ok := o.processes[base.Process_BuilderBackend]; ok {
			t.Fatal("builder process should have been removed from processes map")
		}
		if _, ok := o.processes[base.Process_EditorBackend]; !ok {
			t.Fatal("editor process should not be touched by StopGame()")
		}
		if o.IsGameRunning() {
			t.Fatal("IsGameRunning() = true after StopGame()")
		}
	})

	// Also the cleanup path for failed or crashed runs, so it has to be
	// callable when there is nothing to stop.
	t.Run("no game running", func(t *testing.T) {
		o := NewOrchestrator()
		o.StopGame()

		if len(o.processes) != 0 {
			t.Fatalf("expected empty processes map, got %d", len(o.processes))
		}
	})
}

func Test_Reverse(t *testing.T) {
	input := []string{"a", "b", "c"}
	got := reverse(input)

	want := []string{"c", "b", "a"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("reverse() = %v, want %v", got, want)
		}
	}

	// Shutdown order derives from startup order, so this must not mutate it.
	if input[0] != "a" {
		t.Fatalf("reverse() mutated its input: %v", input)
	}
}

// The editor and the monitor can both decide a run is over at once, and a
// process must only ever be stopped (and waited on) once.
func Test_StopGame_ConcurrentStopsClaimTheProcessOnce(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start test process: %v", err)
	}

	o := NewOrchestrator()
	o.processes[base.Process_BuilderBackend] = &entity.Process{
		Name: base.Process_BuilderBackend,
		Cmd:  cmd,
	}

	const callers = 8
	var wg sync.WaitGroup
	wg.Add(callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			o.StopGame()
		}()
	}
	wg.Wait()

	if o.IsGameRunning() {
		t.Fatal("IsGameRunning() = true after concurrent StopGame() calls")
	}
}

// takeProcess is the claim itself: exactly one caller may win a given process.
func Test_TakeProcess_OnlyOneCallerWins(t *testing.T) {
	o := NewOrchestrator()
	o.processes[base.Process_BuilderBackend] = &entity.Process{Name: base.Process_BuilderBackend}

	const callers = 16
	var wg sync.WaitGroup
	var wins int64
	wg.Add(callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			if _, ok := o.takeProcess(base.Process_BuilderBackend); ok {
				atomic.AddInt64(&wins, 1)
			}
		}()
	}
	wg.Wait()

	if wins != 1 {
		t.Fatalf("takeProcess() succeeded %d times, want exactly 1", wins)
	}
}
