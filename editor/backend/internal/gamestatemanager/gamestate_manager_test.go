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
			gameobject.Gameobject_CurLabelComponents: map[string]interface{}{},
			gameobject.Gameobject_CurLabelID:         obj.GetID(),
			gameobject.Gameobject_CurLabelName:       "",
			gameobject.Gameobject_CurLabelGroup:      "",
		},
	}

	state := gsm.GetGameState()

	assert.Equal(t, expectState, state)
}

func Test_Reset(t *testing.T) {
	gsm := NewGameStateManager()
	gsm.AddGameobject(gameobject.NewGameobject())
	gsm.AddGameobject(gameobject.NewGameobject())

	gsm.Reset()

	assert.Equal(t, 0, len(gsm.GetGameState()))
}

func Test_BuildFromDetails(t *testing.T) {
	t.Run("replaces the existing scene rather than merging", func(t *testing.T) {
		gsm := NewGameStateManager()
		gsm.AddGameobject(gameobject.NewGameobject())

		err := gsm.BuildFromDetails(map[string]interface{}{
			"obj1": map[string]interface{}{
				gameobject.Gameobject_CurLabelName:  "player",
				gameobject.Gameobject_CurLabelGroup: "",
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(gsm.GetGameState()))

		obj, found := gsm.GetGameobject("obj1")
		assert.Equal(t, true, found)
		assert.Equal(t, "player", obj.GetName())
	})

	t.Run("invalid scene shape returns an error and does not partially apply", func(t *testing.T) {
		gsm := NewGameStateManager()

		err := gsm.BuildFromDetails(map[string]interface{}{
			"obj1": "not-a-map",
		})

		assert.Equal(t, ErrExpectedMapStringInterface, err)
		assert.Equal(t, 0, len(gsm.GetGameState()))
	})
}
