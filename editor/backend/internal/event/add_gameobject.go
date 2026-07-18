package event

import (
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// AddGameobjectEvent triggers the creation of new gameobject in the scene
type AddGameobjectEvent struct{}

func NewAddGameobjectEvent() *AddGameobjectEvent {
	return &AddGameobjectEvent{}
}

// Returns the name of the event
func (age *AddGameobjectEvent) GetEventName() string {
	return base.Event_AddGameobject
}

// AddGameobjectEventHandler handles the AddGameobjectEvent and performs the actual functionality
type AddGameobjectEventHandler struct {
	gsm *gamestatemanager.GameStateManager
}

func NewAddGameobjectEventHandler(gsm *gamestatemanager.GameStateManager) *AddGameobjectEventHandler {
	return &AddGameobjectEventHandler{
		gsm: gsm,
	}
}

// update the game scene according to the event
func (agh *AddGameobjectEventHandler) Handle(event models.Event) *EventResponse {
	_, ok := event.(*AddGameobjectEvent)

	if !ok {
		return NewEmittableResponse(false, base.ErrIncorrectEventDataFound, nil)
	}

	gameobject := gameobject.NewGameobject()

	agh.gsm.AddGameobject(gameobject)

	return NewNonEmittableResponse(true, "gameobject created successfully", gameobject.GetGameobjectDetails())
}

var _ models.Event = &AddGameobjectEvent{}
var _ EventHandler = &AddGameobjectEventHandler{}
