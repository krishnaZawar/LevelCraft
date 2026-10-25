package gameobject

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component"
	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/stretchr/testify/assert"
)

const (
	componentName = "component"
)

type MockComponent struct {
	mockGetComponentName    func() string
	mockGetComponentDetails func() component.ComponentDetails
	mockBuildFromDetails    func(json.RawMessage) error
	mockHandlesInput        func(string) bool
}

func (mc *MockComponent) GetComponentName() string {
	return mc.mockGetComponentName()
}
func (mc *MockComponent) GetComponentDetails() component.ComponentDetails {
	return mc.mockGetComponentDetails()
}
func (mc *MockComponent) BuildFromDetails(data json.RawMessage) error {
	return mc.mockBuildFromDetails(data)
}
func (mc *MockComponent) HandlesInput(input string) bool {
	return mc.mockHandlesInput(input)
}

var (
	comp = &MockComponent{
		mockGetComponentName: func() string {
			return componentName
		},
		mockGetComponentDetails: func() component.ComponentDetails {
			return component.ComponentDetails{
				Name: componentName,
				Data: json.RawMessage(`{
					"field1" : "val1",
					"field2": {
						"field: "val"
					}
				}`),
			}
		},
		mockBuildFromDetails: func(data json.RawMessage) error {
			return nil
		},
		mockHandlesInput: func(input string) bool {
			return false
		},
	}

	registry = component.NewComponentRegistry(&input.InputMapping{})
)

func Test_NewGameobject(t *testing.T) {
	obj := NewGameobject(registry)

	assert.Equal(t, DefaultGameobjectName, obj.GetName())
	assert.Equal(t, DefaultGameobjectGroup, obj.GetGroup())

	name, group := "name", "group"

	obj.SetGroup(group)
	obj.SetName(name)

	assert.Equal(t, name, obj.GetName())
	assert.Equal(t, group, obj.GetGroup())
}

func Test_NewGameobjectWithID(t *testing.T) {
	obj := NewGameobjectWithID("123", registry)

	assert.Equal(t, "123", obj.GetID())
	assert.Equal(t, DefaultGameobjectName, obj.GetName())
	assert.Equal(t, DefaultGameobjectGroup, obj.GetGroup())
}

func Test_AddComponent(t *testing.T) {
	obj := NewGameobject(registry)

	t.Run("when component does not exist", func(t *testing.T) {
		ok := obj.AddComponent(comp)
		assert.Equal(t, true, ok)

		val, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, true, ok)
		assert.Equal(t, comp, val)
	})

	t.Run("when component exists", func(t *testing.T) {
		ok := obj.AddComponent(comp)
		assert.Equal(t, false, ok)
	})
}

func Test_RemoveComponent(t *testing.T) {
	obj := NewGameobject(registry)
	obj.AddComponent(comp)
	t.Run("delete non existent component", func(t *testing.T) {
		obj.RemoveComponent("")
		val, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, true, ok)
		assert.Equal(t, comp, val)
	})
	t.Run("delete existing component", func(t *testing.T) {
		obj.RemoveComponent(comp.GetComponentName())
		_, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, false, ok)
	})
}

func Test_GetComponent(t *testing.T) {
	obj := NewGameobject(registry)

	t.Run("when component does not exist", func(t *testing.T) {
		val, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, false, ok)
		assert.Nil(t, val)
	})

	t.Run("when component exists", func(t *testing.T) {
		ok := obj.AddComponent(comp)
		assert.Equal(t, true, ok)

		val, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, true, ok)
		assert.Equal(t, comp, val)
	})
}

func Test_ListComponents(t *testing.T) {
	obj := NewGameobject(registry)

	t.Run("no components", func(t *testing.T) {
		list := obj.ListComponents()
		assert.Equal(t, []component.Component{}, list)
	})

	t.Run("when components present", func(t *testing.T) {
		ok := obj.AddComponent(comp)
		assert.Equal(t, true, ok)

		list := obj.ListComponents()
		assert.Equal(t, []component.Component{comp}, list)
	})
}

func Test_BuildFromDetails(t *testing.T) {
	obj := NewGameobject(registry)
	obj.registry = registry

	id := obj.GetID()
	name, group := "name", "group"

	var comp component.Component = component.NewTransform(100, 100, 100, 100)

	data := GameobjectDetails{
		Id:    obj.GetID(),
		Name:  name,
		Group: group,
		Components: []component.ComponentDetails{
			comp.GetComponentDetails(),
		},
	}

	err := obj.BuildFromDetails(data)
	assert.Nil(t, err)

	objComp, found := obj.GetComponent(base.ComponentName_Transform)

	assert.Equal(t, true, found)
	assert.Equal(t, id, obj.GetID())
	assert.Equal(t, name, obj.GetName())
	assert.Equal(t, group, obj.GetGroup())
	assert.Equal(t, 1, len(obj.components))
	assert.Equal(t, comp.GetComponentDetails(), objComp.GetComponentDetails())
}
func Test_GetGameobjectDetails(t *testing.T) {
	obj := NewGameobject(registry)

	name, group := "name", "group"
	obj.SetGroup(group)
	obj.SetName(name)

	_ = obj.AddComponent(comp)

	data := obj.GetGameobjectDetails()

	/*
		expected result:

		map[
			components:list[
				map[
					field1:val1
					field2:map[field:val]
				]
			]
			group:group
			id:7376d80e-12d1-4465-bdd4-8c788e603e45
			name:name
		]
	*/
	compData := comp.GetComponentDetails()
	expectedData := GameobjectDetails{
		Id:    obj.GetID(),
		Name:  name,
		Group: group,
		Components: []component.ComponentDetails{
			compData,
		},
	}

	assert.Equal(t, data, expectedData)
}
