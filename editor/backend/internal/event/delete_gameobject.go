package event

import (
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// DeleteGameobjectEvent triggers the deletion of gameobject in the scene, if it exists
type DeleteGameobjectEvent struct {
	ID string // unique identifier of the gameobject to be deleted
}

func NewDeleteGameobjectEvent(id string) *DeleteGameobjectEvent {
	return &DeleteGameobjectEvent{
		ID: id,
	}
}

// Returns the name of the event
func (dge *DeleteGameobjectEvent) GetEventName() string {
	return base.Event_DeleteGameobject
}

// DeleteGameobjectHandler handles the DeleteGameobjectEvent and performs the actual functionality
type DeleteGameobjectEventHandler struct {
	gsm *gamestatemanager.GameStateManager
}

func NewDeleteGameobjectEventHandler(gsm *gamestatemanager.GameStateManager) *DeleteGameobjectEventHandler {
	return &DeleteGameobjectEventHandler{
		gsm: gsm,
	}
}

// update the game scene according to the event
func (dgh *DeleteGameobjectEventHandler) Handle(event models.Event) *EventResponse {
	evt, ok := event.(*DeleteGameobjectEvent)
	if !ok {
		return NewEmittableResponse(false, base.ErrIncorrectEventDataFound, nil)
	}

	dgh.gsm.DeleteGameobject(evt.ID)

	return NewNonEmittableResponse(true, "gameobject deleted successfully", nil)
}

var _ models.Event = &DeleteGameobjectEvent{}
var _ EventHandler = &DeleteGameobjectEventHandler{}
