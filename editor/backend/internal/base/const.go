package base

const ServiceName = "game-editor"

// holds the names of all the commands used by the editor
const (
	Command_AddGameobject    = "AddGameobjectCommand"
	Command_DeleteGameobject = "DeleteGameobjectCommand"
	Command_AttachComponent  = "AttachComponentCommand"
	Command_DetachComponent  = "DetachComponentCommand"
	Command_UpdateComponent  = "UpdateComponentCommand"
)

// holds the names of all the events used by the editor
const (
	Event_AddGameobject    = "AddGameobjectEvent"
	Event_DeleteGameobject = "DeleteGameobjectEvent"
	Event_AttachComponent  = "AttachComponentEvent"
	Event_DetachComponent  = "DetachComponentEvent"
	Event_UpdateComponent  = "UpdateComponentEvent"
)

const (
	// max number of events that can be processed in one frame
	MaxEventsProcessablePerFrame = 1000

	// max number of commands that can be processed in one frame
	MaxCommandsProcessablePerFrame = 100
)

const (
	// event assertion failure message
	ErrIncorrectEventDataFound = "incorrect event data found"
)
