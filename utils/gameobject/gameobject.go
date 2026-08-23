package gameobject

import (
	"errors"

	"github.com/google/uuid"
	"github.com/krishnaZawar/LevelCraft/utils/component"
)

var (
	ErrExpectedString               = errors.New("GameobjectError: Expected string value")
	ErrComponentsStructureIncorrect = errors.New("GamobjectError: Components structure was incorrect")
	ErrExpectedMapStringInterface   = errors.New("GameobjectError: Expected map string interface")
	ErrComponentNotFound            = errors.New("GameobjectError: Component not found")
	ErrObjectIDCorrupted            = errors.New("GameobjectError: ObjectID was corrupted")
)

// Used to define the labels used for marshalling and unmarshalling
// defined as []string to provide backward compatibility in future
var (
	gameobject_LabelsID         = []string{"id"}
	gameobject_LabelsName       = []string{"name"}
	gameobject_LabelsGroup      = []string{"group"}
	gameobject_LabelsComponents = []string{"components"}
)

// current value used for unmarshalling
// should be added to the labels slice
const (
	Gameobject_CurLabelID         = "id"
	Gameobject_CurLabelName       = "name"
	Gameobject_CurLabelGroup      = "group"
	Gameobject_CurLabelComponents = "components"
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

func NewGameobject() *Gameobject {
	return &Gameobject{
		id:         uuid.NewString(),
		components: make(map[string]component.Component),
		registry:   component.NewComponentRegistry(),
	}
}

func NewGameobjectWithID(id string) *Gameobject {
	return &Gameobject{
		id:         id,
		components: make(map[string]component.Component),
		registry:   component.NewComponentRegistry(),
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
		Gameobject_CurLabelID:         g.id,
		Gameobject_CurLabelName:       g.name,
		Gameobject_CurLabelGroup:      g.group,
		Gameobject_CurLabelComponents: componentsData,
	}
}

// Build gameobject from provided details
func (g *Gameobject) BuildFromDetails(data map[string]interface{}) error {
	temp := *g
	for _, val := range gameobject_LabelsID {
		if v, ok := data[val]; ok {
			switch n := v.(type) {
			case string:
				if temp.id != n {
					return ErrObjectIDCorrupted
				}
			default:
				return ErrExpectedString
			}
			break
		}
	}

	for _, val := range gameobject_LabelsName {
		if v, ok := data[val]; ok {
			switch n := v.(type) {
			case string:
				temp.name = n
			default:
				return ErrExpectedString
			}
			break
		}
	}

	for _, val := range gameobject_LabelsGroup {
		if v, ok := data[val]; ok {
			switch n := v.(type) {
			case string:
				temp.group = n
			default:
				return ErrExpectedString
			}
			break
		}
	}

	for _, val := range gameobject_LabelsComponents {
		if v, ok := data[val]; ok {
			switch n := v.(type) {
			case map[string]interface{}:
				for compName, data := range n {
					comp, found := temp.registry.GetComponent(compName)
					if !found {
						return ErrComponentNotFound
					}
					compData, ok := data.(map[string]interface{})
					if !ok {
						return ErrComponentsStructureIncorrect
					}
					err := comp.BuildFromDetails(compData)
					if err != nil {
						return err
					}
					temp.components[compName] = comp
				}
			default:
				return ErrExpectedMapStringInterface
			}
			break
		}
	}

	*g = temp

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
