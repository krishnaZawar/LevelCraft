package command

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
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

func Test_KeyDownCommandHandle(t *testing.T) {
	t.Run("KeyDownCommandTest valid data to handle function", func(t *testing.T) {
		comm := &KeyDownCommand{
			KeyName:      Name_KeyA,
			inputManager: inputmanager.NewInputManager(&testMapping),
		}
		val, err := comm.inputManager.GetInputState(comm.KeyName)
		assert.Nil(t, err)
		assert.Equal(t, inputmanager.KeyState_Up, val)

		comm.Handle()

		val, err = comm.inputManager.GetInputState(comm.KeyName)
		assert.Nil(t, err)
		assert.Equal(t, inputmanager.KeyState_Down, val)
	})

	t.Run("KeyDownCommandTest invalid data to handle function", func(t *testing.T) {
		comm := &KeyDownCommand{
			KeyName:      Name_NonExistentKey,
			inputManager: inputmanager.NewInputManager(&testMapping),
		}
		val, err := comm.inputManager.GetInputState(comm.KeyName)
		assert.Equal(t, inputmanager.ErrKeyNotFound, err)
		assert.Equal(t, inputmanager.KeyState_Null, val)

		comm.Handle()

		val, err = comm.inputManager.GetInputState(comm.KeyName)
		assert.Equal(t, inputmanager.ErrKeyNotFound, err)
		assert.Equal(t, inputmanager.KeyState_Null, val)
	})
}

func Test_KeyUpCommandHandle(t *testing.T) {
	t.Run("KeyUpCommandTest valid data to handle function", func(t *testing.T) {
		comm := &KeyUpCommand{
			KeyName:      Name_KeyA,
			inputManager: inputmanager.NewInputManager(&testMapping),
		}
		val, err := comm.inputManager.GetInputState(comm.KeyName)
		assert.Nil(t, err)
		assert.Equal(t, inputmanager.KeyState_Up, val)

		comm.Handle()

		val, err = comm.inputManager.GetInputState(comm.KeyName)
		assert.Nil(t, err)
		assert.Equal(t, inputmanager.KeyState_Up, val)
	})

	t.Run("KeyUpCommandTest invalid data to handle function", func(t *testing.T) {
		comm := &KeyUpCommand{
			KeyName:      Name_NonExistentKey,
			inputManager: inputmanager.NewInputManager(&testMapping),
		}
		val, err := comm.inputManager.GetInputState(comm.KeyName)
		assert.Equal(t, inputmanager.ErrKeyNotFound, err)
		assert.Equal(t, inputmanager.KeyState_Null, val)

		comm.Handle()

		val, err = comm.inputManager.GetInputState(comm.KeyName)
		assert.Equal(t, inputmanager.ErrKeyNotFound, err)
		assert.Equal(t, inputmanager.KeyState_Null, val)
	})
}
