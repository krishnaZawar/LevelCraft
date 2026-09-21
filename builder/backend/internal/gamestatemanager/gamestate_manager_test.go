package gamestatemanager

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component"
	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/stretchr/testify/assert"
)

func Test_AddGameobject(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.AddGameobject(gameobject.NewGameobject())

	assert.Equal(t, 1, len(gsm.GetGameState()))
}

func Test_DeleteGameobject(t *testing.T) {
	gsm := NewGameStateManager()
	obj := gameobject.NewGameobject()

	gsm.AddGameobject(obj)

	t.Run("delete non existing gameobject", func(t *testing.T) {
		gsm.DeleteGameobject("")

		assert.Equal(t, 1, len(gsm.GetGameState()))
	})
	t.Run("delete existing gameobject", func(t *testing.T) {
		gsm.DeleteGameobject(obj.GetID())

		assert.Equal(t, 0, len(gsm.GetGameState()))
	})
}

func Test_GetGameobject(t *testing.T) {
	gsm := NewGameStateManager()
	expectedObj := gameobject.NewGameobject()

	gsm.AddGameobject(expectedObj)

	obj, found := gsm.GetGameobject(expectedObj.GetID())

	assert.Equal(t, expectedObj, obj)
	assert.Equal(t, true, found)
}

func Test_GetGameState(t *testing.T) {
	defVal := 100
	gsm := NewGameStateManager()
	obj := gameobject.NewGameobject()
	obj.AddComponent(component.NewTransform(defVal, defVal, defVal, defVal))

	gsm.AddGameobject(obj)

	expectState := []gameobject.GameobjectDetails{
		{
			Id:    obj.GetID(),
			Name:  obj.GetName(),
			Group: obj.GetGroup(),
			Components: []component.ComponentDetails{
				{
					Name: base.ComponentName_Transform,
					Data: json.RawMessage(`{
						"x": 100,
						"y": 100,
						"w": 100,
						"h": 100
					}`),
				},
			},
		},
	}

	state := gsm.GetGameState()

	assert.Equal(t, expectState[0].Id, state[0].Id)
	assert.Equal(t, expectState[0].Name, state[0].Name)
	assert.Equal(t, expectState[0].Group, state[0].Group)
	assert.Equal(t, expectState[0].Components[0].Name, state[0].Components[0].Name)
	assert.JSONEq(t, string(expectState[0].Components[0].Data), string(state[0].Components[0].Data))
}

func Test_buildFromDetails(t *testing.T) {
	t.Run("valid build", func(t *testing.T) {
		gsm := NewGameStateManager()

		id, name, group := "obj123", "name", "group"

		sceneData := []byte(`[
			{
				"id": "obj123",
				"name": "name",
				"group": "group",
				"components": [
					{
						"name": "Transform",
						"data": {
							"x": 100,
							"y": 100,
							"w": 100,
							"h": 100
						}
					}
				]
			}
		]`)

		compData := component.ComponentDetails{
			Name: base.ComponentName_Transform,
			Data: json.RawMessage(`{
				"x": 100,
				"y": 100,
				"w": 100,
				"h": 100
			}`),
		}

		err := gsm.BuildFromDetails(sceneData)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(gsm.gameobjects))

		obj, found := gsm.GetGameobject(id)
		assert.Equal(t, true, found)
		assert.Equal(t, id, obj.GetID())
		assert.Equal(t, name, obj.GetName())
		assert.Equal(t, group, obj.GetGroup())

		comp, found := obj.GetComponent(base.ComponentName_Transform)
		assert.Equal(t, true, found)

		actualDetails := comp.GetComponentDetails()
		assert.JSONEq(t, string(compData.Data), string(actualDetails.Data))
		assert.Equal(t, compData.Name, actualDetails.Name)
	})
	t.Run("invalid build", func(t *testing.T) {
		gsm := NewGameStateManager()

		sceneData := []byte(`[
			{
				"id": "obj123",
				"name": "name",
				"group": "group",
				"components": [
					{
						"name": "transform",
						"Data": {
							"x": 100,
							"y": 100,
							"w": 100,
							"h": 100
						}
					}
				]
			}
		]`)

		err := gsm.BuildFromDetails(sceneData)
		assert.NotNil(t, err)
		assert.Equal(t, 0, len(gsm.gameobjects))
	})
}
