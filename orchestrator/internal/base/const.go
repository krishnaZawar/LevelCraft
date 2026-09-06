package base

import "time"

const (
	ServiceName = "orchestrator"
)

// holds the process names that will start for the application to run
const (
	Process_EditorBackend  = "editorBackend"
	Process_EditorFrontend = "editorFrontend"
)

const (
	MaxStartAttempts = 3                      // the maximum number of attempts the process takes to start successfully
	PingTimeout      = 2 * time.Second        // the http client timeout for performing the ping function
	PingPollInterval = 300 * time.Millisecond // how often the orchestrator polls /ping while waiting for a process to come up

	MonitorInterval = 5 * time.Second // time interval after which the orchestrator monitors the health of all the services
)

// how long a single start attempt waits for /ping before retrying
const (
	BackendStartupTimeout  = 10 * time.Second
	FrontendStartupTimeout = 45 * time.Second // Electron + Vite cold start is slower
)

const (
	MinPortValue = 1024  // min value of port that can be used by any process
	MaxPortValue = 49151 // max value of port that can be used by any process
)

// env vars passed to the editor frontend process
const (
	EnvEditorBackendURL = "LEVELCRAFT_BACKEND_URL"
	EnvEditorPingPort   = "LEVELCRAFT_PING_PORT"
)
