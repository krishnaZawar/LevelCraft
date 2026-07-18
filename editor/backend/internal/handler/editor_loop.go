package handler

import (
	"time"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/eventmanager"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/queue"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
	utilqueue "github.com/krishnaZawar/LevelCraft/utils/queue"
)

// The editor loop is where all the command and event processing happens
func EditorLoop() {
	cmdQueue := queue.GetCommandQueue()
	evtQueue := queue.GetEventQueue()
	respQueue := queue.GetRespQueue()

	evtManager := eventmanager.Get()

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		updateEditor(cmdQueue, evtQueue, respQueue, evtManager)
		resp := event.EventResponse{
			Success: true,
			Msg:     "gameState",
			Data:    gamestatemanager.Get().GetGameState(),
		}
		respQueue.Push(resp)
	}
}

// Processes and update the state each iteration
func updateEditor(
	cmdQueue *utilqueue.CommandQueue,
	evtQueue *utilqueue.EventQueue,
	respQueue *helper.Queue[event.EventResponse],
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

		ls.Info().Msgf("resp received for %s is %+v", event.GetEventName(), resp)

		if resp.ShouldEmit {
			respQueue.Push(*resp)
		}
		eventsProcessed++
	}

	// perform simulations
}
