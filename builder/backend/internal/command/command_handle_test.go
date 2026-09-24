package command

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	"github.com/stretchr/testify/assert"
)

const (
	Name_KeyA  = "KEY_A"
	Code_KeyA  = 65
	Label_KeyA = "A"

	Name_Digit1  = "DIGIT_1"
	Code_Digit1  = 49
	Label_Digit1 = "1"

	Name_LMB  = "MOUSE_BUTTON_LEFT"
	Code_LMB  = 1
	Label_LMB = "left mouse button"

	Name_RMB  = "MOUSE_BUTTON_RIGHT"
	Code_RMB  = 2
	Label_RMB = "right mouse button"

	Name_NonExistentKey = "NonExistentKey"
)

var (
	input_KeyA = input.InputData{
		Name:  Name_KeyA,
		Code:  Code_KeyA,
		Label: Label_KeyA,
	}
	input_Digit1 = input.InputData{
		Name:  Name_Digit1,
		Code:  Code_Digit1,
		Label: Label_Digit1,
	}

	input_LMB = input.InputData{
		Name:  Name_LMB,
		Code:  Code_LMB,
		Label: Label_LMB,
	}
	input_RMB = input.InputData{
		Name:  Name_RMB,
		Code:  Code_RMB,
		Label: Label_RMB,
	}

	testMapping = input.InputMapping{
		Keyboard: []input.InputData{
			input_KeyA, input_Digit1,
		},
		Mouse: []input.InputData{
			input_LMB, input_RMB,
		},
	}
)

func Test_CommandHandles(t *testing.T) {
	tests := []struct {
		name                        string
		keyName                     string
		createComm                  func(im *inputmanager.InputManager) models.Command
		expectedInitialFetchErr     bool
		expectedFetchErrAfterHandle bool
		expectedInitialKeyState     inputmanager.KeyState
		expectedFinalKeyState       inputmanager.KeyState
	}{
		{
			name:    "KeyDownCommandTest valid data to handle function",
			keyName: Name_KeyA,
			createComm: func(im *inputmanager.InputManager) models.Command {
				return &KeyDownCommand{
					KeyName:      Name_KeyA,
					inputManager: im,
				}
			},
			expectedInitialFetchErr:     false,
			expectedFetchErrAfterHandle: false,
			expectedInitialKeyState:     inputmanager.KeyState_Up,
			expectedFinalKeyState:       inputmanager.KeyState_Down,
		},
		{
			name:    "KeyDownCommandTest invalid data to handle function",
			keyName: Name_NonExistentKey,
			createComm: func(im *inputmanager.InputManager) models.Command {
				return &KeyDownCommand{
					KeyName:      Name_NonExistentKey,
					inputManager: im,
				}
			},
			expectedInitialFetchErr:     true,
			expectedFetchErrAfterHandle: true,
			expectedInitialKeyState:     inputmanager.KeyState_Null,
			expectedFinalKeyState:       inputmanager.KeyState_Null,
		},
		{
			name:    "KeyUpCommandTest valid data to handle function",
			keyName: Name_KeyA,
			createComm: func(im *inputmanager.InputManager) models.Command {
				return &KeyUpCommand{
					KeyName:      Name_KeyA,
					inputManager: im,
				}
			},
			expectedInitialFetchErr:     false,
			expectedFetchErrAfterHandle: false,
			expectedInitialKeyState:     inputmanager.KeyState_Up,
			expectedFinalKeyState:       inputmanager.KeyState_Up,
		},
		{
			name:    "KeyUpCommandTest invalid data to handle function",
			keyName: Name_NonExistentKey,
			createComm: func(im *inputmanager.InputManager) models.Command {
				return &KeyUpCommand{
					KeyName:      Name_NonExistentKey,
					inputManager: im,
				}
			},
			expectedInitialFetchErr:     true,
			expectedFetchErrAfterHandle: true,
			expectedInitialKeyState:     inputmanager.KeyState_Null,
			expectedFinalKeyState:       inputmanager.KeyState_Null,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputManager := inputmanager.NewInputManager(&testMapping)
			comm := tt.createComm(inputManager)

			val, err := inputManager.GetInputState(tt.keyName)
			if tt.expectedInitialFetchErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedInitialKeyState, val)

			comm.Handle()

			val, err = inputManager.GetInputState(tt.keyName)
			if tt.expectedFetchErrAfterHandle {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedFinalKeyState, val)
		})
	}
}
