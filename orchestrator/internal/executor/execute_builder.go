package executor

import (
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/base"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/entity"
	"github.com/krishnaZawar/LevelCraft/orchestrator/internal/layout"
)

// starts the builder backend for the game scene at scenePath
// The scene is a CLI arg and never reloaded: one scene per process
func StartBuilderBackend(scenePath string) (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_BuilderBackend, func() entity.CommandConfig {
		return builderBackendConfig(layout.Get(), getRandomPort(), scenePath)
	}, base.BackendStartupTimeout)
}

// starts the builder frontend, pointing it at the given builder backend
func StartBuilderFrontend(backendURI string) (*entity.Process, error) {
	return startAndWaitHealthy(base.Process_BuilderFrontend, func() entity.CommandConfig {
		port := getRandomPort()
		comm := builderFrontendConfig(layout.Get(), port)
		comm.Env = []string{
			base.EnvBuilderBackendURL + "=" + backendURI,
			base.EnvBuilderPingPort + "=" + port,
		}
		return comm
	}, base.FrontendStartupTimeout)
}

func builderBackendConfig(l *layout.Layout, port string, scenePath string) entity.CommandConfig {
	if l.IsPackaged() {
		return entity.CommandConfig{
			Name: l.Path(l.BuilderBackend),
			Args: []string{"--port", port, "--file-name", scenePath},
			Port: port,
		}
	}
	return entity.CommandConfig{
		Pwd:  "../builder/backend",
		Name: "go",
		Args: []string{"run", "cmd/main.go", "--port", port, "--file-name", scenePath},
		Port: port,
	}
}

func builderFrontendConfig(l *layout.Layout, port string) entity.CommandConfig {
	if l.IsPackaged() {
		return entity.CommandConfig{
			Name: l.Path(l.BuilderApp),
			Port: port,
		}
	}
	return entity.CommandConfig{
		Pwd:  "../builder/app",
		Name: "npm",
		Args: []string{"run", "dev"},
		Port: port,
	}
}
