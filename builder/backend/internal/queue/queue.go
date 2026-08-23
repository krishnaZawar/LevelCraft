package queue

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/command"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/krishnaZawar/LevelCraft/utils/queue"
)

// This is used to create the decoder and register
// all the commandFactories relevant to the types of Commands for decoding CommandRequests
func initCommandFactoryDecoder() *helper.Registry[string, models.CommandFactory] {
	inputmanager := inputmanager.Get()

	decoder := helper.NewRegistry[string, models.CommandFactory]()
	decoder.Register(base.Command_KeyDown, command.NewKeyDownCommandFactory(inputmanager))
	decoder.Register(base.Command_KeyUp, command.NewKeyUpCommandFactory(inputmanager))
	return decoder
}

var cmdQueue = queue.NewCommandQueue(initCommandFactoryDecoder())

var evtQueue = queue.NewEventQueue()

var respQueue = helper.NewQueue[models.EventResponse]()

func GetCommandQueue() *queue.CommandQueue {
	return cmdQueue
}

func GetEventQueue() *queue.EventQueue {
	return evtQueue
}

func GetRespQueue() *helper.Queue[models.EventResponse] {
	return respQueue
}
