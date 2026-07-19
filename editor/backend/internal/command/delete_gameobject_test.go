package command

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteGameobjectCommand(t *testing.T) {
	dgc := &DeleteGameobjectCommand{}

	assert.Equal(t, base.Command_DeleteGameobject, dgc.GetCommandName())

	events := []models.Event{event.NewDeleteGameobjectEvent("")}
	assert.Equal(t, events, dgc.Handle())
}

func Test_DeleteGameobjectCommandFactory(t *testing.T) {
	dgcf := NewDeleteGameobjectCommandFactory()

	t.Run("valid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123"
		}`)
		comm, err := dgcf.NewCommand(data)

		assert.Nil(t, err)
		assert.IsType(t, &DeleteGameobjectCommand{}, comm)
	})

	t.Run("invalid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123"
		`)
		comm, err := dgcf.NewCommand(data)

		assert.NotNil(t, err)
		assert.Nil(t, comm)
	})
}
