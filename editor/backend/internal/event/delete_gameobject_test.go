package event

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteGameobjectEvent(t *testing.T) {
	dge := NewDeleteGameobjectEvent("123")

	assert.Equal(t, "123", dge.ID)
	assert.Equal(t, base.Event_DeleteGameobject, dge.GetEventName())
}

func Test_DeleteGameobjectEventHandler(t *testing.T) {
	gsm := gamestatemanager.NewGameStateManager()
	dgh := NewDeleteGameobjectEventHandler(gsm)

	t.Run("called with wrong event", func(t *testing.T) {
		resp := dgh.Handle(NewAddGameobjectEvent())

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.ErrIncorrectEventDataFound, resp.Msg)
		assert.Nil(t, resp.Data)
	})

	obj := gameobject.NewGameobject()
	gsm.AddGameobject(obj)
	t.Run("called with wrong id", func(t *testing.T) {
		resp := dgh.Handle(NewDeleteGameobjectEvent(""))

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "gameobject deleted successfully", resp.Msg)
		assert.Equal(t, 1, len(gsm.GetGameState()))
	})
	t.Run("called with correct id", func(t *testing.T) {
		resp := dgh.Handle(NewDeleteGameobjectEvent(obj.GetID()))

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "gameobject deleted successfully", resp.Msg)
		assert.Equal(t, 0, len(gsm.GetGameState()))
	})
}
