package event

import (
	"fmt"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// DetachComponentEvent triggers the removal of the component, if exists, from a specific gameobject in the scene
type DetachComponentEvent struct {
	ID            string // unique identifier of the gameobject from which the component is removed
	ComponentName string // name of the component to be removed
}

func NewDetachComponentEvent(id string, name string) *DetachComponentEvent {
	return &DetachComponentEvent{
		ID:            id,
		ComponentName: name,
	}
}

// Returns the name of the event
func (dce *DetachComponentEvent) GetEventName() string {
	return base.Event_DetachComponent
}

// DetachComponentEventHandler handles the DetachComponentEvent and performs the actual functionality
type DetachComponentEventHandler struct {
	gsm *gamestatemanager.GameStateManager
}

func NewDetachComponentEventHandler(gsm *gamestatemanager.GameStateManager) *DetachComponentEventHandler {
	return &DetachComponentEventHandler{
		gsm: gsm,
	}
}

// update the game scene according to the event
func (dch *DetachComponentEventHandler) Handle(event models.Event) *EventResponse {
	evt, ok := event.(*DetachComponentEvent)
	if !ok {
		return NewEmittableResponse(false, base.ErrIncorrectEventDataFound, nil)
	}

	obj, found := dch.gsm.GetGameobject(evt.ID)
	if !found {
		return NewNonEmittableResponse(false, fmt.Sprintf("gameobject with id %s not found", evt.ID), nil)
	}

	obj.RemoveComponent(evt.ComponentName)

	return NewNonEmittableResponse(true, "component removed successfully", nil)
}

var _ models.Event = &DetachComponentEvent{}
var _ EventHandler = &DetachComponentEventHandler{}
