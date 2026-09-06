package executor

import (
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

// starts the editor frontend, pointing it at the given editor backend
func StartEditorFrontend(backendURI string) (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_EditorFrontend, func() entity.CommandConfig {
		port := getRandomPort()
		return entity.CommandConfig{
			Pwd:  "../editor/app",
			Name: "npm",
			Args: []string{
				"run", "dev",
			},
			Port: port,
			Env: []string{
				base.EnvEditorBackendURL + "=" + backendURI,
				base.EnvEditorPingPort + "=" + port,
			},
		}
	}, base.FrontendStartupTimeout)
}
