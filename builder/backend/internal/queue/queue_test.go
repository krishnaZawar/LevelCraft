package queue

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/command"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_InitCommandFactoryDecoder(t *testing.T) {
	inputmanager := &inputmanager.InputManager{}
	decoder := initCommandFactoryDecoder(inputmanager)

	assert.Equal(t, 2, decoder.Size())

	tests := []struct {
		key      string
		expected models.CommandFactory
	}{
		{
			key:      base.Command_KeyDown,
			expected: command.NewKeyDownCommandFactory(inputmanager),
		},
		{
			key:      base.Command_KeyUp,
			expected: command.NewKeyUpCommandFactory(inputmanager),
		},
	}

	for _, tt := range tests {
		cf, found := decoder.GetValue(tt.key)
		assert.Equal(t, true, found)
		assert.Equal(t, tt.expected, cf)
	}
}
