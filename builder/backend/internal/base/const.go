package base

const ServiceName = "game-builder"

const (
	Command_KeyDown = "KeyDownCommand"
	Command_KeyUp   = "KeyUpCommand"
)

const (
	// max number of events that can be processed in one frame
	MaxEventsProcessablePerFrame = 1000

	// max number of commands that can be processed in one frame
	MaxCommandsProcessablePerFrame = 100
)
