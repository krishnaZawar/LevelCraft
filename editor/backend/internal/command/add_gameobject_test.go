package command

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_AddGameobjectCommand(t *testing.T) {
	agc := &AddGameobjectCommand{}

	assert.Equal(t, base.Command_AddGameobject, agc.GetCommandName())

	events := []models.Event{event.NewAddGameobjectEvent()}
	assert.Equal(t, events, agc.Handle())
}

func Test_AddGameobjectCommandFactory(t *testing.T) {
	agcf := NewAddGameobjectCommandFactory()
	comm, err := agcf.NewCommand([]byte{})
	assert.IsType(t, &AddGameobjectCommand{}, comm)
	assert.Nil(t, err)
}
