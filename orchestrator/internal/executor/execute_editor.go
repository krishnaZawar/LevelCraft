package executor

import (
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

// starts the editor backend
func StartEditorBackend() (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_EditorBackend, func() entity.CommandConfig {
		port := getRandomPort()
		return entity.CommandConfig{
			Pwd:  "../editor/backend",
			Name: "go",
			Args: []string{
				"run", "cmd/main.go", "--port", port,
			},
			Port: port,
		}
	}, base.BackendStartupTimeout)
}
