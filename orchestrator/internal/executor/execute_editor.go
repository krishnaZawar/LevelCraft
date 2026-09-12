package executor

import (
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/layout"
)

// starts the editor backend
func StartEditorBackend() (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_EditorBackend, func() entity.CommandConfig {
		return editorBackendConfig(layout.Get(), getRandomPort())
	}, base.BackendStartupTimeout)
}

// starts the editor frontend, pointing it at the editor backend and at the
// orchestrator's control API, which it uses to ask for game runs
func StartEditorFrontend(backendURI string, orchestratorURI string) (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_EditorFrontend, func() entity.CommandConfig {
		port := getRandomPort()
		comm := editorFrontendConfig(layout.Get(), port)
		comm.Env = []string{
			base.EnvEditorBackendURL + "=" + backendURI,
			base.EnvEditorPingPort + "=" + port,
			base.EnvOrchestratorURL + "=" + orchestratorURI,
		}
		return comm
	}, base.FrontendStartupTimeout)
}

func editorBackendConfig(l *layout.Layout, port string) entity.CommandConfig {
	if l.IsPackaged() {
		return entity.CommandConfig{
			Name: l.Path(l.EditorBackend),
			Args: []string{"--port", port},
			Port: port,
		}
	}
	return entity.CommandConfig{
		Pwd:  "../editor/backend",
		Name: "go",
		Args: []string{"run", "cmd/main.go", "--port", port},
		Port: port,
	}
}

func editorFrontendConfig(l *layout.Layout, port string) entity.CommandConfig {
	if l.IsPackaged() {
		return entity.CommandConfig{
			Name: l.Path(l.EditorApp),
			Port: port,
		}
	}
	return entity.CommandConfig{
		Pwd:  "../editor/app",
		Name: "npm",
		Args: []string{"run", "dev"},
		Port: port,
	}
}
