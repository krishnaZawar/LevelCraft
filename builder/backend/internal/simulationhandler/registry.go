package simulationhandler

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

// builds the registry to hold the simulation handlers corresponding to the components
func initSimulationHandlerRegistry(gsm *gamestatemanager.GameStateManager) *helper.Registry[string, SimulationHandler] {
	registry := helper.NewRegistry[string, SimulationHandler]()

	registry.Register(compbase.ComponentName_ObjectFollow, NewObjectFollowSimulator(gsm))

	return registry
}

var eventGeneratorRegistry = initSimulationHandlerRegistry(gamestatemanager.Get())

// returns the singleton registry that will be used across the engine for fetching simulationHandlers
func GetSimulationHandlerRegistry() *helper.Registry[string, SimulationHandler] {
	return eventGeneratorRegistry
}
