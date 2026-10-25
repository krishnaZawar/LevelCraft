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

	// Length of the frame being generated for, in seconds.
	//
	// Frame generators need it to size the last step of a move: without it an
	// object steps past the point it was heading for and spends every frame
	// correcting, which reads as a jitter around that point.
	DeltaTime float64
}

// Handles the generation and events specific to the component
type EventGenerator interface {
	// Generate events based on the component type and the input received
	GenerateEvents(EventGeneratorCtx) ([]models.Event, error)
}
