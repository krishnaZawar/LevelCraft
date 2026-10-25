package gameobject

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/krishnaZawar/LevelCraft/utils/component"
)

var (
	ErrExpectedString               = errors.New("GameobjectError: Expected string value")
	ErrComponentsStructureIncorrect = errors.New("GamobjectError: Components structure was incorrect")
	ErrObjectIDCorrupted            = errors.New("GameobjectError: ObjectID was corrupted")

	NewErrComponentNotFound = func(compName string) error {
		return fmt.Errorf("GameobjectError: Component %s not found", compName)
	}
)

const (
	DefaultGameobjectName  = "Gameobject"
	DefaultGameobjectGroup = "Group"
)

// Gameobject is the container that represents an object in the scene.
//
// The main function of the Gameobject is to organize and hold Components of an object together
type Gameobject struct {
	id         string                         // unique identifier of the gameobject
	name       string                         // name of the gameobject
	group      string                         // the group that the gameobject belongs to
	components map[string]component.Component // collection of all the Components held by the gameobject

	registry *component.ComponentRegistry // injected to build the gameobject from details provided
}

// structure of Gameobject used for marshalling and unmarshalling gameobject details and holding object details
// the actual update tasks are done on the Gameobject struct only
type GameobjectDetails struct {
	Id         string                       `json:"id"`
	Name       string                       `json:"name"`
	Group      string                       `json:"group"`
	Components []component.ComponentDetails `json:"components"`
}

func NewGameobject(registry *component.ComponentRegistry) *Gameobject {
	obj := &Gameobject{
		id:         uuid.NewString(),
		components: make(map[string]component.Component),
		registry:   registry,
	}
	obj.SetName(DefaultGameobjectName)
	obj.SetGroup(DefaultGameobjectGroup)
	return obj
}

func NewGameobjectWithID(id string, registry *component.ComponentRegistry) *Gameobject {
	obj := &Gameobject{
		id:         id,
		components: make(map[string]component.Component),
		registry:   registry,
	}
	obj.SetName(DefaultGameobjectName)
	obj.SetGroup(DefaultGameobjectGroup)
	return obj
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

// List all the components the gameobject holds
func (g *Gameobject) ListComponents() []component.Component {
	comps := []component.Component{}
	for _, comp := range g.components {
		comps = append(comps, comp)
	}
	return comps
}

// Returns the all the details of the gameobject
func (g *Gameobject) GetGameobjectDetails() GameobjectDetails {
	componentsData := []component.ComponentDetails{}
	for componentName := range g.components {
		componentsData = append(componentsData, g.components[componentName].GetComponentDetails())
	}
	return GameobjectDetails{
		Id:         g.GetID(),
		Name:       g.GetName(),
		Group:      g.GetGroup(),
		Components: componentsData,
	}
}

// Build gameobject from provided details
func (g *Gameobject) BuildFromDetails(data GameobjectDetails) error {
	if g.id != data.Id {
		return ErrObjectIDCorrupted
	}

	g.SetName(data.Name)
	g.SetGroup(data.Group)

	for _, compDetails := range data.Components {
		factory, found := g.registry.GetComponent(compDetails.Name)
		if !found {
			return NewErrComponentNotFound(compDetails.Name)
		}
		comp := factory()
		err := comp.BuildFromDetails(compDetails.Data)
		if err != nil {
			return err
		}
		g.AddComponent(comp)
	}

	return nil
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
