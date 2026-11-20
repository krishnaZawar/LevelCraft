package command

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/controlmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/input"
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

func Test_KeyCommandHandles(t *testing.T) {
	tests := []struct {
		name                        string
		keyName                     string
		createComm                  func(im *inputmanager.InputManager) models.Command
		expectedInitialFetchErr     bool
		expectedFetchErrAfterHandle bool
		expectedHandleErr           bool
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
			expectedHandleErr:           false,
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
			expectedHandleErr:           true,
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
			expectedHandleErr:           false,
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
			expectedHandleErr:           true,
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

			err = comm.Handle()
			if tt.expectedHandleErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

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

func Test_ExecutionStateCommandHandles(t *testing.T) {
	defaultExecutionState := controlmanager.ExecutionState_Inactive
	tests := []struct {
		name                        string
		createComm                  func(cm *controlmanager.ControlManager) models.Command
		expectedHandleErr           bool
		expectedFinalExecutionState controlmanager.ExecutionState
	}{
		{
			name: "SetExecutionStateCommandTest valid data to handle function",
			createComm: func(cm *controlmanager.ControlManager) models.Command {
				return &SetExecutionStateCommand{
					ExecutionState: controlmanager.ExecutionState_Active,
					controlManager: cm,
				}
			},
			expectedHandleErr:           false,
			expectedFinalExecutionState: controlmanager.ExecutionState_Active,
		},
		{
			name: "SetExecutionStateCommandTest invalid data to handle function",
			createComm: func(cm *controlmanager.ControlManager) models.Command {
				return &SetExecutionStateCommand{
					ExecutionState: controlmanager.ExecutionState("invalid type"),
					controlManager: cm,
				}
			},
			expectedHandleErr:           true,
			expectedFinalExecutionState: defaultExecutionState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &controlmanager.ControlManager{}
			_ = cm.SetExecutionState(defaultExecutionState)

			comm := tt.createComm(cm)

			err := comm.Handle()
			if tt.expectedHandleErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.expectedFinalExecutionState, cm.GetExecutionState())
		})
	}
}
