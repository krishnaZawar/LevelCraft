package command

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_CommandNames(t *testing.T) {
	tests := []struct {
		command      models.Command
		expectedName string
	}{
		{
			command:      &KeyDownCommand{},
			expectedName: base.Command_KeyDown,
		},
		{
			command:      &KeyUpCommand{},
			expectedName: base.Command_KeyUp,
		},
	}

	for _, tt := range tests {
		name := tt.command.GetCommandName()
		assert.Equal(t, tt.expectedName, name)
	}
}
