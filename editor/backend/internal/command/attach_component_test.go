package command

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_AttachComponentCommand(t *testing.T) {
	acc := &AttachComponentCommand{}

	assert.Equal(t, base.Command_AttachComponent, acc.GetCommandName())

	events := []models.Event{event.NewAttachComponentEvent("", "")}
	assert.Equal(t, events, acc.Handle())
}

func Test_AttachComponentCommandFactory(t *testing.T) {
	accf := NewAttachComponentCommandFactory()
	t.Run("valid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123",
			"componentName": "ok"
		}`)
		comm, err := accf.NewCommand(data)

		assert.Nil(t, err)
		assert.IsType(t, &AttachComponentCommand{}, comm)
	})

	t.Run("invalid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123"
		`)
		comm, err := accf.NewCommand(data)

		assert.NotNil(t, err)
		assert.Nil(t, comm)
	})
}
