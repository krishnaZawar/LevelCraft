package gamestatemanager

import (
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

func Test_buildFromDetails(t *testing.T) {
	t.Run("valid build", func(t *testing.T) {
		gsm := NewGameStateManager()

		id := "obj123"
		name, group := "name", "group"
		defVal := 100

		compDetails := map[string]interface{}{
			component.Transform_CurLabelX: defVal,
			component.Transform_CurLabelY: defVal,
			component.Transform_CurLabelW: defVal,
			component.Transform_CurLabelH: defVal,
		}

		buildData := map[string]interface{}{
			id: map[string]interface{}{
				gameobject.Gameobject_CurLabelID:    id,
				gameobject.Gameobject_CurLabelName:  name,
				gameobject.Gameobject_CurLabelGroup: group,
				gameobject.Gameobject_CurLabelComponents: map[string]interface{}{
					base.ComponentName_Transform: compDetails,
				},
			},
		}

		err := gsm.BuildFromDetails(buildData)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(gsm.gameobjects))

		obj, found := gsm.GetGameobject(id)
		assert.Equal(t, true, found)
		assert.Equal(t, id, obj.GetID())
		assert.Equal(t, name, obj.GetName())
		assert.Equal(t, group, obj.GetGroup())

		comp, found := obj.GetComponent(base.ComponentName_Transform)
		assert.Equal(t, true, found)
		assert.Equal(t, compDetails, comp.GetComponentDetails())
	})
	t.Run("invalid build", func(t *testing.T) {
		gsm := NewGameStateManager()

		id := "obj123"
		name, group := "name", "group"
		defVal := 100

		compDetails := map[string]interface{}{
			component.Transform_CurLabelX: defVal,
			component.Transform_CurLabelY: defVal,
			component.Transform_CurLabelW: defVal,
			component.Transform_CurLabelH: defVal,
		}

		buildData := map[string]interface{}{
			id: map[string]interface{}{
				gameobject.Gameobject_CurLabelID:    id + "1",
				gameobject.Gameobject_CurLabelName:  name,
				gameobject.Gameobject_CurLabelGroup: group,
				gameobject.Gameobject_CurLabelComponents: map[string]interface{}{
					base.ComponentName_Transform: compDetails,
				},
			},
		}

		err := gsm.BuildFromDetails(buildData)
		assert.NotNil(t, err)
		assert.Equal(t, 0, len(gsm.gameobjects))
	})
}
