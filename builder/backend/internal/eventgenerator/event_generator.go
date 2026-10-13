package eventgenerator

import (
	"github.com/krishnaZawar/LevelCraft/utils/component"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// Holds the details required for the event generation from context
type EventGeneratorCtx struct {
	GameobjectID string              // gameobject id linked with the interacting component
	Component    component.Component // interacting component object
	InputKey     string              // input key with which the component can interact
}

// Handles the generation and events specific to the component
type EventGenerator interface {
	// Generate events based on the component type and the input received
	GenerateEvents(EventGeneratorCtx) ([]models.Event, error)
}
