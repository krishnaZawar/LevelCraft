package executor

import (
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
)

// starts the builder backend for the game scene at scenePath
// The scene is a CLI arg and never reloaded: one scene per process
func StartBuilderBackend(scenePath string) (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_BuilderBackend, func() entity.CommandConfig {
		return builderBackendConfig(getRandomPort(), scenePath)
	}, base.BackendStartupTimeout)
}

// starts the builder frontend, pointing it at the given builder backend
func StartBuilderFrontend(backendURI string) (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_BuilderFrontend, func() entity.CommandConfig {
		port := getRandomPort()
		comm := builderFrontendConfig(port)
		comm.Env = []string{
			base.EnvBuilderBackendURL + "=" + backendURI,
			base.EnvBuilderPingPort + "=" + port,
		}
		return comm
	}, base.FrontendStartupTimeout)
}

func builderBackendConfig(port string, scenePath string) entity.CommandConfig {
	return entity.CommandConfig{
		Pwd:  "../builder/backend",
		Name: "go",
		Args: []string{"run", "cmd/main.go", "--port", port, "--file-name", scenePath},
		Port: port,
	}
}

func builderFrontendConfig(port string) entity.CommandConfig {
	return entity.CommandConfig{
		Pwd:  "../builder/app",
		Name: "npm",
		Args: []string{"run", "dev"},
		Port: port,
	}
}
