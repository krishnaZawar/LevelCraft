package event

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/gamestatemanager"
	"github.com/stretchr/testify/assert"
)

func Test_AddGameobjectEvent(t *testing.T) {
	age := NewAddGameobjectEvent()

	assert.Equal(t, base.Event_AddGameobject, age.GetEventName())
}

func Test_AddGameobjectEventHandler(t *testing.T) {
	gsm := gamestatemanager.NewGameStateManager()
	agh := NewAddGameobjectEventHandler(gsm)

	t.Run("called with wrong event", func(t *testing.T) {
		resp := agh.Handle(NewDeleteGameobjectEvent("123"))

		assert.Equal(t, false, resp.Success)
		assert.Equal(t, base.ErrIncorrectEventDataFound, resp.Msg)
		assert.Nil(t, resp.Data)
	})

	t.Run("called with correct event", func(t *testing.T) {
		resp := agh.Handle(NewAddGameobjectEvent())

		assert.Equal(t, true, resp.Success)
		assert.Equal(t, "gameobject created successfully", resp.Msg)
		assert.Equal(t, 1, len(gsm.GetGameState()))
	})
}
