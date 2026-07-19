package queue

import (
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/command"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/krishnaZawar/LevelCraft/utils/queue"
)

// This is used to create the decoder and register
// all the commandFactories relevant to the types of Commands for decoding CommandRequests
func initCommandFactoryDecoder() *helper.Registry[string, models.CommandFactory] {
	decoder := helper.NewRegistry[string, models.CommandFactory]()
	decoder.Register(base.Command_AddGameobject, command.NewAddGameobjectCommandFactory())
	decoder.Register(base.Command_DeleteGameobject, command.NewDeleteGameobjectCommandFactory())
	decoder.Register(base.Command_AttachComponent, command.NewAttachComponentCommandFactory())
	decoder.Register(base.Command_DetachComponent, command.NewDetachComponentCommandFactory())
	decoder.Register(base.Command_UpdateComponent, command.NewUpdateComponentCommandFactory())
	return decoder
}

var cmdQueue = queue.NewCommandQueue(initCommandFactoryDecoder())

var evtQueue = queue.NewEventQueue()

var respQueue = helper.NewQueue[event.EventResponse]()

func GetCommandQueue() *queue.CommandQueue {
	return cmdQueue
}

func GetEventQueue() *queue.EventQueue {
	return evtQueue
}

func GetRespQueue() *helper.Queue[event.EventResponse] {
	return respQueue
}
