package event

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	componentbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/stretchr/testify/assert"
)

func Test_DetachComponentEvent(t *testing.T) {
	dce := NewDetachComponentEvent("123", "name")

	assert.Equal(t, "123", dce.ID)
	assert.Equal(t, "name", dce.ComponentName)
	assert.Equal(t, base.Event_DetachComponent, dce.GetEventName())
}

func Test_DetachComponentEventHandler(t *testing.T) {
	gsm := gamestatemanager.NewGameStateManager()
	dch := NewDetachComponentEventHandler(gsm)

	t.Run("called with wrong event", func(t *testing.T) {
		resp := dch.Handle(NewAddGameobjectEvent())

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.ErrIncorrectEventDataFound, resp.Msg)
		assert.Nil(t, resp.Data)
	})

	obj := gameobject.NewGameobject()
	gsm.AddGameobject(obj)
	t.Run("called with wrong object id", func(t *testing.T) {
		resp := dch.Handle(NewDetachComponentEvent("1", ""))

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, "gameobject with id 1 not found", resp.Msg)
	})
	comp := component.NewTransform(100, 100, 100, 100)
	obj.AddComponent(comp)
	t.Run("called with wrong component name", func(t *testing.T) {
		resp := dch.Handle(NewDetachComponentEvent(obj.GetID(), "123"))

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "component removed successfully", resp.Msg)
		assert.Equal(t, 1, len(obj.GetGameobjectDetails()["components"].(map[string]interface{})))
	})
	t.Run("called with right object id and component name", func(t *testing.T) {
		resp := dch.Handle(NewDetachComponentEvent(obj.GetID(), componentbase.ComponentName_Transform))

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "component removed successfully", resp.Msg)
		assert.Equal(t, 0, len(obj.GetGameobjectDetails()["components"].(map[string]interface{})))
	})
}
