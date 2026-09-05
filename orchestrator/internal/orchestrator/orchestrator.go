package orchestrator

import (
	"context"
	"time"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/executor"
	"github.com/krishnaZawar/LevelCraft/utils/logger"
)

var ls = logger.New(base.ServiceName)

// handles the type of state the application holds during monitoring
type applicationState string

const (
	applicationState_healthy     applicationState = "healthy"     // no issues in the application
	applicationState_builderExit applicationState = "builderExit" // one of the builder processes exit
	applicationState_editorExit  applicationState = "editorExit"  // one of the editor processes exit
)

// Orchestrator handles the application lifecycle and lifecycle of all its processes
type Orchestrator struct {
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
	err := o.runProcess(executor.StartEditorBackend)
	if err != nil {
		return err
	}
	return nil
}

// monitors all the processes are working fine via healthChecks
func (o *Orchestrator) monitor() applicationState {
	for name, process := range o.processes {
		err := executor.CheckProcessHealth(process.CommunicationURI)
		if err != nil {
			ls.ErrorWith(err).Msgf("failed to get response from process %s", name)
			switch name {
			case base.Process_EditorBackend:
				return applicationState_editorExit
			}
		}
	}
	return applicationState_healthy
}

// functional abstraction to call the execution function for process creation
func (o *Orchestrator) runProcess(fn func() (*entity.Process, error)) error {
	process, err := fn()
	if err != nil {
		return err
	}
	o.processes[process.Name] = process
	return nil
}

// exits the application in the intended order or graceful process cleanup
func (o *Orchestrator) exitApplication() {
	// holds the order of how the cleanup should occur, from first to last
	processes := []string{
		base.Process_EditorBackend,
	}
	o.exitProcesses(processes)

	ls.Info().Msg("application exited")
}

// exits all the processNames processes in the given order
func (o *Orchestrator) exitProcesses(processNames []string) {
	for _, name := range processNames {
		process, ok := o.processes[name]
		if !ok {
			ls.Warn().Msgf("process with name %s does not exist", name)
			continue
		}
		ls.Info().Msgf("Stopping %s", name)

		if err := process.Stop(); err != nil {
			ls.ErrorWith(err).Msgf("failed to exit from process %s", name)
			delete(o.processes, name)
			continue
		}
	}
}
