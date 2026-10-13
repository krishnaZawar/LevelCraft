package handler

import (
	"time"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/eventgenerator"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/eventmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/queue"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// The game loop is where all the command and event processing happens
func GameLoop() {
	cmdQueue := queue.GetCommandQueue()
	evtQueue := queue.GetEventQueue()

	evtManager := eventmanager.Get()
	inputManager := inputmanager.Get()
	gamestateManager := gamestatemanager.Get()

	evtGeneratorRegistry := eventgenerator.GetEventGeneratorRegistry()

	const FPS = 60
	frameStep := 1.0 / FPS

	curFrameStart := time.Now()
	prevFrameStart := curFrameStart

	running := true

	for running {
		curFrameStart = time.Now()
		dt := curFrameStart.Sub(prevFrameStart).Seconds()
		prevFrameStart = curFrameStart

		// process commands and update states
		commandsProcessed := 0
		for cmdQueue.Length() > 0 && commandsProcessed < base.MaxCommandsProcessablePerFrame {
			err := cmdQueue.ConsumeCommand()
			if err != nil {
				ls.ErrorWith(err).Msg("command consumption failed")
				continue
			}
			commandsProcessed++
		}

		// parse through all input states and the game scene and emit relevant events
		gameobjects := gamestateManager.ListGameobjects()
		for input, state := range inputManager.GetAllStates() {
			if state != inputmanager.KeyState_Down {
				continue
			}
			for _, gameobject := range gameobjects {
				for _, comp := range gameobject.ListComponents() {
					if comp.HandlesInput(input) {
						generator, found := evtGeneratorRegistry.GetValue(comp.GetComponentName())
						if !found {
							ls.Warn().Msgf("generator for component %s not found", comp.GetComponentName())
							continue
						}
						eventCtx := eventgenerator.EventGeneratorCtx{
							GameobjectID: gameobject.GetID(),
							Component:    comp,
							InputKey:     input,
						}
						events, err := generator.GenerateEvents(eventCtx)
						if err != nil {
							ls.ErrorWith(err).Msgf("failed to generate events from component %s for ctx %+v", comp.GetComponentName(), eventCtx)
							continue
						}
						for _, event := range events {
							evtQueue.Ingest(event)
						}
					}
				}
			}
		}

		// process the events
		eventsProcessed := 0
		for evtQueue.Length() > 0 && eventsProcessed < base.MaxEventsProcessablePerFrame {
			event, err := evtQueue.ConsumeEvent()
			if err != nil {
				continue
			}
			handlers, found := evtManager.GetHandler(event.GetEventName())
			if !found {
				ls.Error().Msgf("failed to get handler for event %s", event.GetEventName())
				continue
			}
			for _, handler := range handlers {
				resp := handler.Handle(event, dt)
				if !resp.Success {
					ls.Warn().Msgf("event failed, event: %+v, resp: %+v", event, resp)
				}
			}
			eventsProcessed++
		}

		// perform simulations
		resp := models.EventResponse{
			Success: true,
			Msg:     "gamestate",
			Data:    gamestatemanager.Get().GetGameState(),
		}
		queue.GetRespQueue().Push(resp)

		if rem := time.Since(curFrameStart).Seconds(); rem < frameStep {
			time.Sleep(time.Duration((frameStep - rem) * float64(time.Second)))
		}
	}
}
