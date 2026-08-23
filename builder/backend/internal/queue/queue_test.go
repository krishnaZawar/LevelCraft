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
	decoder := initCommandFactoryDecoder()

	tests := []struct {
		key          string
		expectedType models.CommandFactory
	}{
		{
			key:          base.Command_KeyDown,
			expectedType: command.NewKeyDownCommandFactory(&inputmanager.InputManager{}),
		},
	}

	for _, tt := range tests {
		cf, found := decoder.GetValue(tt.key)
		assert.Equal(t, true, found)
		assert.IsType(t, tt.expectedType, cf)
	}
}
