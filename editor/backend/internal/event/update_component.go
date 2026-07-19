package event

import (
	"fmt"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// UpdateComponentEvent triggers the updation of the component, if exists, from a specific gameobject in the scene
type UpdateComponentEvent struct {
	ID            string                 // unique identifier of the gameobject from which the component is updated
	ComponentName string                 // name of the component to be updated
	Data          map[string]interface{} // data to update the component
}

func NewUpdateComponentEvent(id string, name string, data map[string]interface{}) *UpdateComponentEvent {
	return &UpdateComponentEvent{
		ID:            id,
		ComponentName: name,
		Data:          data,
	}
}

// Returns the name of the event
func (dce *UpdateComponentEvent) GetEventName() string {
	return base.Event_UpdateComponent
}

// UpdateComponentEventHandler handles the UpdateComponentEvent and performs the actual functionality
type UpdateComponentEventHandler struct {
	gsm *gamestatemanager.GameStateManager
}

func NewUpdateComponentEventHandler(gsm *gamestatemanager.GameStateManager) *UpdateComponentEventHandler {
	return &UpdateComponentEventHandler{
		gsm: gsm,
	}
}

// update the game scene according to the event
func (uch *UpdateComponentEventHandler) Handle(event models.Event) *EventResponse {
	evt, ok := event.(*UpdateComponentEvent)
	if !ok {
		return NewEmittableResponse(false, base.ErrIncorrectEventDataFound, nil)
	}

	obj, found := uch.gsm.GetGameobject(evt.ID)
	if !found {
		return NewNonEmittableResponse(false, fmt.Sprintf("gameobject with id %s not found", evt.ID), nil)
	}

	comp, found := obj.GetComponent(evt.ComponentName)
	if !found {
		return NewNonEmittableResponse(false, fmt.Sprintf("%s component not found for update", evt.ComponentName), nil)
	}

	comp.BuildFromDetails(evt.Data)

	return NewNonEmittableResponse(true, "component updated successfully", nil)
}

var _ models.Event = &UpdateComponentEvent{}
var _ EventHandler = &UpdateComponentEventHandler{}
