package queue

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/command"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/controlmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_InitCommandFactoryDecoder(t *testing.T) {
	inputmanager := &inputmanager.InputManager{}
	controlManager := &controlmanager.ControlManager{}
	decoder := initCommandFactoryDecoder(inputmanager, controlManager)

	assert.Equal(t, 3, decoder.Size())

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
		{
			key:      base.Command_SetExecutionState,
			expected: command.NewSetExecutionStateCommandFactory(controlManager),
		},
	}

	for _, tt := range tests {
		cf, found := decoder.GetValue(tt.key)
		assert.Equal(t, true, found)
		assert.Equal(t, tt.expected, cf)
	}
}
