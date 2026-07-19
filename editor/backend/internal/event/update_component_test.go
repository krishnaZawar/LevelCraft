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

func Test_UpdateComponentEvent(t *testing.T) {
	uce := NewUpdateComponentEvent("123", "name", map[string]interface{}{})

	assert.Equal(t, "123", uce.ID)
	assert.Equal(t, "name", uce.ComponentName)
	assert.Equal(t, map[string]interface{}{}, uce.Data)
	assert.Equal(t, base.Event_UpdateComponent, uce.GetEventName())
}

func Test_UpdateComponentEventHandler(t *testing.T) {
	gsm := gamestatemanager.NewGameStateManager()
	uch := NewUpdateComponentEventHandler(gsm)

	t.Run("called with wrong event", func(t *testing.T) {
		resp := uch.Handle(NewAddGameobjectEvent())

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.ErrIncorrectEventDataFound, resp.Msg)
		assert.Nil(t, resp.Data)
	})

	obj := gameobject.NewGameobject()
	gsm.AddGameobject(obj)
	t.Run("called with wrong object id", func(t *testing.T) {
		resp := uch.Handle(NewUpdateComponentEvent("1", "", map[string]interface{}{}))

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, "gameobject with id 1 not found", resp.Msg)
	})
	comp := component.NewTransform(100, 100, 100, 100)
	obj.AddComponent(comp)
	t.Run("called with wrong component name", func(t *testing.T) {
		resp := uch.Handle(NewUpdateComponentEvent(obj.GetID(), "123", map[string]interface{}{}))

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, "123 component not found for update", resp.Msg)
	})
	t.Run("called with right object id and component name", func(t *testing.T) {
		resp := uch.Handle(NewUpdateComponentEvent(obj.GetID(), componentbase.ComponentName_Transform, map[string]interface{}{"x": 10}))

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "component updated successfully", resp.Msg)
		val, ok := obj.GetComponent(componentbase.ComponentName_Transform)
		assert.Equal(t, true, ok)

		x, y := val.(*component.Transform).GetPosition()
		w, h := val.(*component.Transform).GetDimension()
		assert.Equal(t, 10, x)
		assert.Equal(t, 100, y)
		assert.Equal(t, 100, w)
		assert.Equal(t, 100, h)
	})
}
