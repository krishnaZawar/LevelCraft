package handler

import (
	"time"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/collisionsystem"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
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

	collisionEvtQueue := queue.GetCollisionEventQueue()
	collisionSystem := collisionsystem.Get()

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

		// detect collisions
		gameobjects = gamestateManager.ListGameobjects()
		for outerIdx := 0; outerIdx < len(gameobjects); outerIdx++ {
			for innerIdx := outerIdx + 1; innerIdx < len(gameobjects); innerIdx++ {
				if collisionSystem.DetectCollision(gameobjects[innerIdx], gameobjects[outerIdx]) {
					/*
						Ingest the event in both the queues so that the resolution can be done on priority without hindrance of maxProcessing from a separate queue.
						While the main eventQueue can listen for the CollisionOccuredEvent and trigger the respective component handlers.
						This will enable collision resolution as priority while still unblocking the events that would be triggered from this Event
					*/
					collisionEvtQueue.Ingest(&event.CollisionOccuredEvent{
						Obj1_ID: gameobjects[innerIdx].GetID(),
						Obj2_ID: gameobjects[outerIdx].GetID(),
					})
					// evtQueue.Ingest(&event.CollisionOccuredEvent{
					// 	Obj1_ID: gameobjects[innerIdx].GetID(),
					// 	Obj2_ID: gameobjects[outerIdx].GetID(),
					// })
				}
			}
		}

		// resolve collisions
		for collisionEvtQueue.Length() > 0 {
			event, err := collisionEvtQueue.ConsumeEvent()
			if err != nil {
				continue
			}
			resp := collisionSystem.Handle(event, dt)
			if !resp.Success {
				ls.Warn().Msgf("collision resolution failed, event: %+v, resp: %+v", event, resp)
			}
		}

		// process pending events

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
