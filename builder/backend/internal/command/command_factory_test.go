package command

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

func Test_CommandFactory(t *testing.T) {
	inputManager := &inputmanager.InputManager{}
	tests := []struct {
		name          string
		factory       models.CommandFactory
		data          json.RawMessage
		ExpectedErr   bool
		ExpectedValue models.Command
	}{
		{
			name:    "KeyDownCommandFactoryTest valid marshal",
			factory: NewKeyDownCommandFactory(inputManager),
			data: json.RawMessage(`{
				"keyName": "KEY_A"
			}`),
			ExpectedErr: false,
			ExpectedValue: &KeyDownCommand{
				inputManager: inputManager,
				KeyName:      "KEY_A",
			},
		},
		{
			name:    "KeyDownCommandFactoryTest valid marshal wrong payload",
			factory: NewKeyDownCommandFactory(inputManager),
			data: json.RawMessage(`{
				"keyNam": "KEY_A"
			}`),
			ExpectedErr: false,
			ExpectedValue: &KeyDownCommand{
				inputManager: inputManager,
			},
		},
		{
			name:    "KeyDownCommandFactoryTest invalid marshal",
			factory: NewKeyDownCommandFactory(inputManager),
			data: json.RawMessage(`{
				"keyNam": "KEY_A
			}`),
			ExpectedErr:   true,
			ExpectedValue: nil,
		},
		{
			name:    "KeyUpCommandFactoryTest valid marshal",
			factory: NewKeyUpCommandFactory(inputManager),
			data: json.RawMessage(`{
				"keyName": "KEY_A"
			}`),
			ExpectedErr: false,
			ExpectedValue: &KeyUpCommand{
				inputManager: inputManager,
				KeyName:      "KEY_A",
			},
		},
		{
			name:    "KeyUpCommandFactoryTest valid marshal wrong payload",
			factory: NewKeyUpCommandFactory(inputManager),
			data: json.RawMessage(`{
				"keyNam": "KEY_A"
			}`),
			ExpectedErr: false,
			ExpectedValue: &KeyUpCommand{
				inputManager: inputManager,
			},
		},
		{
			name:    "KeyUpCommandFactoryTest invalid marshal",
			factory: NewKeyUpCommandFactory(inputManager),
			data: json.RawMessage(`{
				"keyNam": "KEY_A
			}`),
			ExpectedErr:   true,
			ExpectedValue: nil,
		},
	}

	for _, tt := range tests {
		comm, err := tt.factory.NewCommand(tt.data)
		if tt.ExpectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tt.ExpectedValue, comm)
	}
}
