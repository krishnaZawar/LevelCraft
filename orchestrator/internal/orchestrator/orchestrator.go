package orchestrator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/executor"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/server"
	"github.com/krishnaZawar/LevelCraft/utils/logger"
)

var ls = logger.New(base.ServiceName)

var (
	ErrGameAlreadyRunning = errors.New("error: a game is already running")
)

// handles the type of state the application holds during monitoring
type applicationState string

const (
	applicationState_healthy     applicationState = "healthy"     // no issues in the application
	applicationState_builderExit applicationState = "builderExit" // one of the builder processes exit
	applicationState_editorExit  applicationState = "editorExit"  // one of the editor processes exit
)

// the processes that make up a game run, in startup order
// Reversed for shutdown, so a frontend stops before the backend it talks to
var builderProcesses = []string{
	base.Process_BuilderBackend,
	base.Process_BuilderFrontend,
}

// the processes that make up the editor, in the order they are started
var editorProcesses = []string{
	base.Process_EditorBackend,
	base.Process_EditorFrontend,
}

// Orchestrator handles the application lifecycle and lifecycle of all its processes
// Every process is spawned and reaped here, so this stays the single record
type Orchestrator struct {
	mu sync.Mutex
	// process name -> process. Guarded by mu: the monitor loop reads it while
	// the control API's handlers mutate it from the HTTP server's goroutines
	processes map[string]*entity.Process
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		processes: make(map[string]*entity.Process),
	}
}

// Runs the orchestrator
// handles the complete lifecycle from start -> Monitoring -> Exit
func (o *Orchestrator) Run(ctx context.Context) error {
	/*
		Hold lifecycles:
			- Start in order
			- Runtime
			- Failure
			- Exit in order
	*/
	err := o.start()
	if err != nil {
		ls.ErrorWith(err).Msg("application startup failed, terminating...")
		o.exitApplication()
		return err
	}

	ticker := time.NewTicker(base.MonitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			switch o.monitor() {
			case applicationState_builderExit:
				// A game run ending is not an application failure
				ls.Info().Msg("builder processes exited, cleaning up the game run")
				o.StopGame()

			case applicationState_editorExit:
				o.exitApplication()
				return nil
			}

		case <-ctx.Done():
			ls.Info().Msg("shutdown signal received")
			o.exitApplication()
			return nil
		}
	}
}

// starts the application in order
func (o *Orchestrator) start() error {
	// comes up first so its address can be handed to the editor at startup
	orchestratorURI, err := server.Start(o)
	if err != nil {
		return err
	}
	ls.Info().Msgf("orchestrator control API listening on %s", orchestratorURI)

	if err := o.runProcess(executor.StartEditorBackend); err != nil {
		return err
	}

	editorBackend, ok := o.getProcess(base.Process_EditorBackend)
	if !ok {
		return errors.New("editor backend process not found")
	}
	if err := o.runProcess(func() (*entity.Process, error) {
		return executor.StartEditorFrontend(editorBackend.CommunicationURI, orchestratorURI)
	}); err != nil {
		return err
	}

	return nil
}

// RunGame starts the builder processes for the scene at scenePath
// One game at a time: a second request is rejected, not silently swapped in
func (o *Orchestrator) RunGame(scenePath string) error {
	if o.IsGameRunning() {
		return ErrGameAlreadyRunning
	}

	if err := o.runProcess(func() (*entity.Process, error) {
		return executor.StartBuilderBackend(scenePath)
	}); err != nil {
		o.StopGame()
		return err
	}

	builderBackend, ok := o.getProcess(base.Process_BuilderBackend)
	if !ok {
		o.StopGame()
		return errors.New("builder backend process not found")
	}

	if err := o.runProcess(func() (*entity.Process, error) {
		return executor.StartBuilderFrontend(builderBackend.CommunicationURI)
	}); err != nil {
		// tear the half-started run down so the next one starts clean
		o.StopGame()
		return err
	}

	ls.Info().Msgf("game run started for scene %s", scenePath)
	return nil
}

// StopGame stops every process belonging to the current game run
// Safe with nothing running, so it also serves as the failed-run cleanup path
func (o *Orchestrator) StopGame() {
	o.exitProcesses(reverse(builderProcesses))
}

// IsGameRunning reports whether any builder process is currently tracked
func (o *Orchestrator) IsGameRunning() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, name := range builderProcesses {
		if _, ok := o.processes[name]; ok {
			return true
		}
	}
	return false
}

// monitors all the processes are working fine via healthChecks
// Editor failures outrank builder ones: the application is going down anyway
func (o *Orchestrator) monitor() applicationState {
	state := applicationState_healthy

	for name, uri := range o.processURIs() {
		err := executor.CheckProcessHealth(uri)
		if err == nil {
			continue
		}
		ls.ErrorWith(err).Msgf("failed to get response from process %s", name)

		switch name {
		case base.Process_EditorBackend, base.Process_EditorFrontend:
			return applicationState_editorExit
		case base.Process_BuilderBackend, base.Process_BuilderFrontend:
			state = applicationState_builderExit
		}
	}

	return state
}

// snapshot of the tracked processes' health check addresses, copied so the
// slow health checks don't hold the lock against the control API
func (o *Orchestrator) processURIs() map[string]string {
	o.mu.Lock()
	defer o.mu.Unlock()
	uris := make(map[string]string, len(o.processes))
	for name, process := range o.processes {
		uris[name] = process.CommunicationURI
	}
	return uris
}

// fetches a tracked process by name
func (o *Orchestrator) getProcess(name string) (*entity.Process, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	process, ok := o.processes[name]
	return process, ok
}

// functional abstraction to call the execution function for process creation
func (o *Orchestrator) runProcess(fn func() (*entity.Process, error)) error {
	process, err := fn()
	if err != nil {
		return err
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.processes[process.Name] = process
	return nil
}

// exits the application in the intended order or graceful process cleanup
func (o *Orchestrator) exitApplication() {
	// torn down first: started last, and depended on by nothing else
	o.StopGame()
	// stop order: frontend before backend, since frontend depends on it
	o.exitProcesses(reverse(editorProcesses))

	ls.Info().Msg("application exited")
}

// exits all the processNames processes in the given order
func (o *Orchestrator) exitProcesses(processNames []string) {
	for _, name := range processNames {
		process, ok := o.takeProcess(name)
		if !ok {
			continue
		}
		ls.Info().Msgf("Stopping %s", name)

		if err := process.Stop(); err != nil {
			ls.ErrorWith(err).Msgf("failed to exit from process %s", name)
		}
	}
}

// removes a process from tracking and returns it, in one atomic step. Claiming
// before stopping keeps the editor's stop and the monitor's from overlapping
func (o *Orchestrator) takeProcess(name string) (*entity.Process, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	process, ok := o.processes[name]
	if ok {
		delete(o.processes, name)
	}
	return process, ok
}

// returns a reversed copy of the given slice
func reverse(names []string) []string {
	reversed := make([]string, len(names))
	for i, name := range names {
		reversed[len(names)-1-i] = name
	}
	return reversed
}

var _ server.GameController = &Orchestrator{}
