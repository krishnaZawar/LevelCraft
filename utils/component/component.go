package component

import (
	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
)

// Component is a modular unit of data and functions attached to the object
//
// It models the behaviour of the object in the scene
type Component interface {
	// Returns the name of the component
	GetComponentName() string

	// Returns a snapshot of the complete data stored in the component
	//
	// Helpful in recursively building the game scene
	GetComponentDetails() map[string]interface{}

	// Builds the component from the data it is provided
	BuildFromDetails(map[string]interface{})
}

// Creates a new component based on the component name
// returns a copy of the base component
type ComponentRegistry struct {
	// name -> component mapping
	registry *helper.Registry[string, Component]
}

func newComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		registry: helper.NewRegistry[string, Component](),
	}
}

// contains the registry object registered with all the components base copies
func NewComponentRegistry() *ComponentRegistry {
	compRegistry := newComponentRegistry()
	compRegistry.register(base.ComponentName_Transform, newBaseTransform())

	return compRegistry
}

// registers a new component with the registry
func (cr *ComponentRegistry) register(name string, comp Component) {
	cr.registry.Register(name, comp)
}

// fetches the base component for the name
func (cr *ComponentRegistry) GetComponent(name string) (Component, bool) {
	comp, ok := cr.registry.GetValue(name)
	return comp, ok
}
