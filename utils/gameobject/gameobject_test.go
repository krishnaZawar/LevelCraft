package gameobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	componentName = "component"
)

type MockComponent struct {
	mockGetComponentName    func() string
	mockGetComponentDetails func() map[string]interface{}
	mockBuildFromDetails    func(map[string]interface{})
}

func (mc *MockComponent) GetComponentName() string {
	return mc.mockGetComponentName()
}
func (mc *MockComponent) GetComponentDetails() map[string]interface{} {
	return mc.mockGetComponentDetails()
}
func (mc *MockComponent) BuildFromDetails(data map[string]interface{}) {
	mc.mockBuildFromDetails(data)
}

var (
	comp = &MockComponent{
		mockGetComponentName: func() string {
			return componentName
		},
		mockGetComponentDetails: func() map[string]interface{} {
			return map[string]interface{}{
				"field1": "val1",
				"field2": map[string]interface{}{
					"field": "val",
				},
			}
		},
		mockBuildFromDetails: func(m map[string]interface{}) {
			// pass
		},
	}
)

func Test_NewGameobject(t *testing.T) {
	obj := NewGameobject()

	name, group := "name", "group"

	obj.SetGroup(group)
	obj.SetName(name)

	assert.Equal(t, name, obj.GetName())
	assert.Equal(t, group, obj.GetGroup())
}

func Test_AddComponent(t *testing.T) {
	obj := NewGameobject()

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
	obj := NewGameobject()
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
	obj := NewGameobject()

	t.Run("when component does not exist", func(t *testing.T) {
		val, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, false, ok)
		assert.Nil(t, val)
	})

	t.Run("when component exists", func(t *testing.T) {
		ok := obj.AddComponent(comp)

		val, ok := obj.GetComponent(comp.GetComponentName())
		assert.Equal(t, true, ok)
		assert.Equal(t, comp, val)
	})
}

func Test_GetGameobjectDetails(t *testing.T) {
	obj := NewGameobject()

	name, group := "name", "group"
	obj.SetGroup(group)
	obj.SetName(name)

	_ = obj.AddComponent(comp)

	data := obj.GetGameobjectDetails()

	/*
		expected result:

		map[
			components:map[
				component:map[
					field1:val1
					field2:map[field:val]
				]
			]
			group:group
			id:7376d80e-12d1-4465-bdd4-8c788e603e45
			name:name
		]
	*/
	compName := comp.GetComponentName()
	compData := comp.GetComponentDetails()
	expectedData := map[string]interface{}{
		"name":  name,
		"group": group,
		"id":    obj.GetID(),
		"components": map[string]interface{}{
			compName: compData,
		},
	}

	assert.Equal(t, data, expectedData)
}
