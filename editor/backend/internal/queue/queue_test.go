package queue

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/command"
	"github.com/stretchr/testify/assert"
)

func Test_InitCommandFactoryDecoder(t *testing.T) {
	decoder := initCommandFactoryDecoder()

	handler, ok := decoder.GetValue(base.Command_AddGameobject)
	assert.Equal(t, true, ok)
	assert.IsType(t, &command.AddGameobjectCommandFactory{}, handler)

	handler, ok = decoder.GetValue(base.Command_DeleteGameobject)
	assert.Equal(t, true, ok)
	assert.IsType(t, &command.DeleteGameobjectCommandFactory{}, handler)
}
