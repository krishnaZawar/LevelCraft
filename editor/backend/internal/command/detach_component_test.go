package command

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_DetachComponentCommand(t *testing.T) {
	acc := &DetachComponentCommand{}

	assert.Equal(t, base.Command_DetachComponent, acc.GetCommandName())

	events := []models.Event{event.NewDetachComponentEvent("", "")}
	assert.Equal(t, events, acc.Handle())
}

func Test_DetachComponentCommandFactory(t *testing.T) {
	dccf := NewDetachComponentCommandFactory()
	t.Run("valid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123",
			"componentName": "ok"
		}`)
		comm, err := dccf.NewCommand(data)

		assert.Nil(t, err)
		assert.IsType(t, &DetachComponentCommand{}, comm)
	})

	t.Run("invalid unmarshal", func(t *testing.T) {
		data := json.RawMessage(`{
			"id": "123"
		`)
		comm, err := dccf.NewCommand(data)

		assert.NotNil(t, err)
		assert.Nil(t, comm)
	})
}
