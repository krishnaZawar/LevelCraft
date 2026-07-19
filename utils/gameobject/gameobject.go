package gameobject

import (
	"github.com/google/uuid"
	"github.com/krishnaZawar/LevelCraft/utils/component"
)

// Gameobject is the container that represents an object in the scene.
//
// The main function of the Gameobject is to organize and hold Components of an object together
type Gameobject struct {
	id         string                         // unique identifier of the gameobject
	name       string                         // name of the gameobject
	group      string                         // the group that the gameobject belongs to
	components map[string]component.Component // collection of all the Components held by the gameobject
}

func NewGameobject() *Gameobject {
	return &Gameobject{
		id:         uuid.NewString(),
		components: make(map[string]component.Component),
	}
}

// Adds a new Component to the gameobject
//
// Rejects the addition if already a Component of that type resides with the gameobject
func (g *Gameobject) AddComponent(comp component.Component) bool {
	_, ok := g.components[comp.GetComponentName()]
	if ok {
		return false
	}
	g.components[comp.GetComponentName()] = comp
	return true
}

// removes an existing component from the gameobject
func (g *Gameobject) RemoveComponent(name string) {
	delete(g.components, name)
}

// Returns the requested component
//
// Return values:
//   - component: the actual component if found, else nil
//   - bool: true if component found, else false
func (g *Gameobject) GetComponent(componentName string) (component.Component, bool) {
	comp, ok := g.components[componentName]
	return comp, ok
}

// Returns the all the details of the gameobject
func (g *Gameobject) GetGameobjectDetails() map[string]interface{} {
	componentsData := map[string]interface{}{}
	for componentName := range g.components {
		componentsData[componentName] = g.components[componentName].GetComponentDetails()
	}
	return map[string]interface{}{
		"id":         g.id,
		"name":       g.name,
		"group":      g.group,
		"components": componentsData,
	}
}

// Updates the group of the gameobject
func (g *Gameobject) SetGroup(group string) {
	g.group = group
}

// Returns the group the gameobject belongs to
func (g *Gameobject) GetGroup() string {
	return g.group
}

// Updates the name of the gameobject
func (g *Gameobject) SetName(name string) {
	g.name = name
}

// Returns the name of the gameobject
func (g *Gameobject) GetName() string {
	return g.name
}

// returns the ID of the gameobject
func (g *Gameobject) GetID() string {
	return g.id
}
