package queue

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/command"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/controlmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

// This is used to create the decoder and register
// all the commandFactories relevant to the types of Commands for decoding CommandRequests
func initCommandFactoryDecoder(
	inputmanager *inputmanager.InputManager,
	controlManager *controlmanager.ControlManager,
) *helper.Registry[string, models.CommandFactory] {
	decoder := helper.NewRegistry[string, models.CommandFactory]()
	decoder.Register(base.Command_KeyDown, command.NewKeyDownCommandFactory(inputmanager))
	decoder.Register(base.Command_KeyUp, command.NewKeyUpCommandFactory(inputmanager))
	decoder.Register(base.Command_SetExecutionState, command.NewSetExecutionStateCommandFactory(controlManager))
	return decoder
}

var cmdQueue = NewCommandQueue(
	initCommandFactoryDecoder(
		inputmanager.Get(),
		controlmanager.Get(),
	),
)

var evtQueue = NewEventQueue()

// The hard cap on processing of events in a frame can sometimes cause restriction of completing the resolution events hence we should have the separate event queue for Collision resolution events
// This is separate from the standard eventQueue because collisions should be resolved for each frame completely else they would lag and defer in the next frame.
var collisionEvtQueue = NewEventQueue()

var respQueue = helper.NewQueue[models.EventResponse]()

func GetCommandQueue() *CommandQueue {
	return cmdQueue
}

func GetEventQueue() *EventQueue {
	return evtQueue
}

func GetRespQueue() *helper.Queue[models.EventResponse] {
	return respQueue
}

func GetCollisionEventQueue() *EventQueue {
	return collisionEvtQueue
}
