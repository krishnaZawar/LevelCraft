package gamestatemanager

import (
	"testing"

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
	gsm := NewGameStateManager()
	obj := gameobject.NewGameobject()

	gsm.AddGameobject(obj)

	/*
		expected state:
		map[
			<id>: map[
				components:map[]
				group:""
				id:<id>
				name:""
			]
		]
	*/
	expectState := map[string]interface{}{
		obj.GetID(): map[string]interface{}{
			"components": map[string]interface{}{},
			"id":         obj.GetID(),
			"name":       "",
			"group":      "",
		},
	}

	state := gsm.GetGameState()

	assert.Equal(t, expectState, state)
}
