package executor

import (
	"fmt"
	"time"

	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

// starts the editor backend
// handles retries if binding fails or process startup fails
// returns the process details to the orchestrator on successful startup
func StartEditorBackend() (*entity.Process, error) {
	var process *entity.Process
	var err error
	attempts := 1
	for attempts <= base.MaxStartAttempts {
		ls.Info().Msgf("starting with attempt %d of starting editor backend", attempts)

		port := getRandomPort()
		comm := entity.CommandConfig{
			Pwd:  "../editor/backend",
			Name: "go",
			Args: []string{
				"run", "cmd/main.go", "--port", port,
			},
			Port: port,
		}

		process, err = buildAndRunProcess(base.Process_EditorBackend, &comm)
		if err != nil {
			ls.ErrorWith(err).Msg("failed to start editor backend")
			attempts++
			continue
		}

		healthy := true
		time.Sleep(base.PingWaitTimeOnStartup)
		err = CheckProcessHealth(process.CommunicationURI)
		if err != nil {
			ls.ErrorWith(err).Msg("failed to ping editor backend. Retrying...")
			time.Sleep(base.PingWaitTimeOnStartup)
			if err := process.Stop(); err != nil {
				healthy = false
				ls.ErrorWith(err).Msgf("failed to ping editor backend on retry. Stopping process: %d", process.Cmd.Process.Pid)
			}
		}
		if healthy {
			break
		}
		if err := process.Stop(); err != nil {
			ls.ErrorWith(err).Msgf("failed to kill process with PID: %d", process.Cmd.Process.Pid)
		}
		attempts++
	}

	if attempts > base.MaxStartAttempts {
		return nil, fmt.Errorf("failed to start editor backend after %d attempts", base.MaxStartAttempts)
	}

	ls.Info().Msgf("Editor Backend started... PID: %d", process.Cmd.Process.Pid)
	return process, nil
}
