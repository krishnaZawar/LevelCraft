package handler

import (
	"time"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/eventmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/queue"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	utilqueue "github.com/krishnaZawar/LevelCraft/utils/queue"
)

// The game loop is where all the command and event processing happens
func GameLoop() {
	cmdQueue := queue.GetCommandQueue()
	evtQueue := queue.GetEventQueue()

	evtManager := eventmanager.Get()

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		updateGame(cmdQueue, evtQueue, evtManager)
		resp := models.EventResponse{
			Success: true,
			Msg:     "inputState",
			Data:    gamestatemanager.Get().GetGameState(),
		}
		queue.GetRespQueue().Push(resp)
	}
}

// Processes and update the state each iteration
func updateGame(
	cmdQueue *utilqueue.CommandQueue,
	evtQueue *utilqueue.EventQueue,
	evtManager *eventmanager.EventManager,
) {
	// process commands and emit out relevant events
	commandsProcessed := 0
	for cmdQueue.Length() > 0 && commandsProcessed < base.MaxCommandsProcessablePerFrame {
		events, err := cmdQueue.ConsumeCommand()
		if err != nil {
			ls.ErrorWith(err).Msg("command consumption failed")
			continue
		}
		for _, event := range events {
			evtQueue.Ingest(event)
		}
		commandsProcessed++
	}
	// process the events
	eventsProcessed := 0
	for evtQueue.Length() > 0 && eventsProcessed < base.MaxEventsProcessablePerFrame {
		event, err := evtQueue.ConsumeEvent()
		if err != nil {
			continue
		}
		handler, found := evtManager.GetHandler(event.GetEventName())
		if !found {
			ls.Error().Msgf("failed to get handler for event %s", event.GetEventName())
			continue
		}
		resp := handler.Handle(event)

		if !resp.Success {
			ls.Warn().Msgf("event failed, event: %+v, resp: %+v", event, resp)
		}

		eventsProcessed++
	}

	// perform simulations
}
