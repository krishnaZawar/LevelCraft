package eventgenerator

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

// builds the registry to hold the event generators corresponding to the components
func initEventGeneratorRegistry(gsm *gamestatemanager.GameStateManager) *helper.Registry[string, EventGenerator] {
	registry := helper.NewRegistry[string, EventGenerator]()

	registry.Register(compbase.ComponentName_InputMovement, NewInputMovementEventGenerator(gsm))

	return registry
}

var eventGeneratorRegistry = initEventGeneratorRegistry(gamestatemanager.Get())

// returns the singleton registry that will be used across the engine for fetching eventGenerators
func GetEventGeneratorRegistry() *helper.Registry[string, EventGenerator] {
	return eventGeneratorRegistry
}
