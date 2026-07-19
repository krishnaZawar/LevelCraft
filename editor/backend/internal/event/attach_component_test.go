package event

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	componentbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/stretchr/testify/assert"
)

func Test_AttachComponentEvent(t *testing.T) {
	dge := NewAttachComponentEvent("123", "name")

	assert.Equal(t, "123", dge.ID)
	assert.Equal(t, "name", dge.ComponentName)
	assert.Equal(t, base.Event_AttachComponent, dge.GetEventName())
}

func Test_AttachComponentEventHandler(t *testing.T) {
	gsm := gamestatemanager.NewGameStateManager()
	dgh := NewAttachComponentEventHandler(gsm)

	t.Run("called with wrong event", func(t *testing.T) {
		resp := dgh.Handle(NewAddGameobjectEvent())

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.ErrIncorrectEventDataFound, resp.Msg)
		assert.Nil(t, resp.Data)
	})

	obj := gameobject.NewGameobject()
	gsm.AddGameobject(obj)
	t.Run("called with wrong object id", func(t *testing.T) {
		resp := dgh.Handle(NewAttachComponentEvent("1", ""))

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, "gameobject with id 1 not found", resp.Msg)
	})
	t.Run("called with wrong component name", func(t *testing.T) {
		resp := dgh.Handle(NewAttachComponentEvent(obj.GetID(), "123"))

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, "could not find base component for 123", resp.Msg)
	})
	t.Run("called with right object id and component name", func(t *testing.T) {
		resp := dgh.Handle(NewAttachComponentEvent(obj.GetID(), componentbase.ComponentName_Transform))

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "component added successfully", resp.Msg)
		assert.Equal(t, 1, len(obj.GetGameobjectDetails()["components"].(map[string]interface{})))
	})
}
