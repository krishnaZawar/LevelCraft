package simulationhandler

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/component"
)

// Holds the details required for the event generation from context
type SimulatorCtx struct {
	GameobjectID string              // gameobject id linked with the simulating component
	Component    component.Component // simulating component object
	Dt           float64             // delta time
}

// Handles the generation and events specific to the component
type SimulationHandler interface {
	// runs continuous simulations for that component based on the current game state
	// optionally emits events for cross component interactions
	Simulate(SimulatorCtx) ([]models.Event, error)
}
