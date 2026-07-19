package event

import (
	"fmt"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// AttachComponentEvent triggers the addition of the component, if not exists to a specific gameobject in the scene
type AttachComponentEvent struct {
	ID            string // unique identifier of the gameobject to which the component is added
	ComponentName string // name of the component to be added
}

func NewAttachComponentEvent(id string, name string) *AttachComponentEvent {
	return &AttachComponentEvent{
		ID:            id,
		ComponentName: name,
	}
}

// Returns the name of the event
func (ace *AttachComponentEvent) GetEventName() string {
	return base.Event_AttachComponent
}

// AttachComponentEventHandler handles the AttachComponentEvent and performs the actual functionality
type AttachComponentEventHandler struct {
	gsm          *gamestatemanager.GameStateManager
	compRegistry *component.ComponentRegistry
}

func NewAttachComponentEventHandler(gsm *gamestatemanager.GameStateManager) *AttachComponentEventHandler {
	return &AttachComponentEventHandler{
		gsm:          gsm,
		compRegistry: component.NewComponentRegistry(),
	}
}

// update the game scene according to the event
func (ach *AttachComponentEventHandler) Handle(event models.Event) *EventResponse {
	evt, ok := event.(*AttachComponentEvent)
	if !ok {
		return NewEmittableResponse(false, base.ErrIncorrectEventDataFound, nil)
	}

	obj, found := ach.gsm.GetGameobject(evt.ID)
	if !found {
		return NewNonEmittableResponse(false, fmt.Sprintf("gameobject with id %s not found", evt.ID), nil)
	}

	comp, ok := ach.compRegistry.GetComponent(evt.ComponentName)
	if !ok {
		return NewNonEmittableResponse(false, fmt.Sprintf("could not find base component for %s", evt.ComponentName), nil)
	}

	ok = obj.AddComponent(comp)
	if !ok {
		return NewNonEmittableResponse(false, fmt.Sprintf("could not add component %s to object %s", evt.ComponentName, evt.ID), nil)
	}

	return NewNonEmittableResponse(true, "component added successfully", nil)
}

var _ models.Event = &AttachComponentEvent{}
var _ EventHandler = &AttachComponentEventHandler{}
