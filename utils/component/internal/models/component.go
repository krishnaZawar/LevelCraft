package models

import (
	"encoding/json"

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
	GetComponentDetails() ComponentDetails

	// Builds the component from the data it is provided
	BuildFromDetails(json.RawMessage) error

	// defines whether a component can interact with given user input or not
	HandlesInput(string) bool
}

// holds the component details structure to type-safe storing of components with heterogeneous attributes
type ComponentDetails struct {
	Name string          `json:"name"` // holds the name of the component
	Data json.RawMessage `json:"data"` // holds the marshalled data of the component
}

// the factory function that is used to create and serve the component
type ComponentFactory func() Component

// Creates a new component based on the component name
// returns a copy of the base component
type ComponentRegistry struct {
	// name -> component mapping
	registry *helper.Registry[string, ComponentFactory]
}

func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		registry: helper.NewRegistry[string, ComponentFactory](),
	}
}

// registers a new component with the registry
func (cr *ComponentRegistry) Register(name string, compFactory ComponentFactory) {
	cr.registry.Register(name, compFactory)
}

// fetches the base component for the name
func (cr *ComponentRegistry) GetComponent(name string) (ComponentFactory, bool) {
	comp, ok := cr.registry.GetValue(name)
	return comp, ok
}

// returns the size of the mapping of the registry
func (cr *ComponentRegistry) Size() int {
	return cr.registry.Size()
}
