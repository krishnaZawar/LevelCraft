package eventgenerator

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

// Generators driven by input: the loop reaches these only while a key the
// component responds to is held down.
func initEventGeneratorRegistry(gsm *gamestatemanager.GameStateManager) *helper.Registry[string, EventGenerator] {
	registry := helper.NewRegistry[string, EventGenerator]()

	registry.Register(compbase.ComponentName_InputMovement, NewInputMovementEventGenerator(gsm))

	return registry
}

// Generators driven by the frame: the loop runs these once per frame for
// every object carrying the component, with no input involved.
//
// A component producing movement on its own has no key press to hang off, so
// without this registry it could never emit anything at all.
func initFrameEventGeneratorRegistry(gsm *gamestatemanager.GameStateManager) *helper.Registry[string, EventGenerator] {
	registry := helper.NewRegistry[string, EventGenerator]()

	registry.Register(compbase.ComponentName_Patrol, NewPatrolEventGenerator(gsm))
	registry.Register(compbase.ComponentName_Follow, NewFollowEventGenerator(gsm))

	return registry
}

var eventGeneratorRegistry = initEventGeneratorRegistry(gamestatemanager.Get())
var frameEventGeneratorRegistry = initFrameEventGeneratorRegistry(gamestatemanager.Get())

// returns the singleton registry that will be used across the engine for fetching eventGenerators
func GetEventGeneratorRegistry() *helper.Registry[string, EventGenerator] {
	return eventGeneratorRegistry
}

// returns the singleton registry of generators that run every frame
func GetFrameEventGeneratorRegistry() *helper.Registry[string, EventGenerator] {
	return frameEventGeneratorRegistry
}
