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
	MaxStartAttempts      = 3               // the maximum number of attempts the process takes to start successfully
	PingWaitTimeOnStartup = 2 * time.Second // the wait time before the orchestrator can ping the process on startup
	PingTimeout           = 2 * time.Second // the http client timeout for performing the ping function

	MonitorInterval = 5 * time.Second // time interval after which the orchestrator monitors the health of all the services
)

const (
	MinPortValue = 1024  // min value of port that can be used by any process
	MaxPortValue = 49151 // max value of port that can be used by any process
)
