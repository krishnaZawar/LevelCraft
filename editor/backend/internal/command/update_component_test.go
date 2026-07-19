package command

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateComponentCommand(t *testing.T) {
	acc := &UpdateComponentCommand{
		Data: map[string]interface{}{},
	}

	assert.Equal(t, base.Command_UpdateComponent, acc.GetCommandName())

	events := []models.Event{event.NewUpdateComponentEvent("", "", map[string]interface{}{})}
	assert.Equal(t, events, acc.Handle())
}

func Test_UpdateComponentCommandFactory(t *testing.T) {
	uccf := NewUpdateComponentCommandFactory()
	t.Run("valid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123",
			"componentName": "ok",
			"data" : {}
		}`)
		comm, err := uccf.NewCommand(data)

		assert.Nil(t, err)
		assert.IsType(t, &UpdateComponentCommand{}, comm)
	})

	t.Run("invalid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123"
		`)
		comm, err := uccf.NewCommand(data)

		assert.NotNil(t, err)
		assert.Nil(t, comm)
	})
}
