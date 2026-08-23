package inputmanager

import (
	"testing"

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
)

var (
	testMapping = input.InputMapping{
		Keyboard: []input.InputData{
			input_KeyA, input_Digit1,
		},
		Mouse: []input.InputData{
			input_LMB, input_RMB,
		},
	}

	expectedTransformedMapping = map[string]input.InputData{
		Name_KeyA:   input_KeyA,
		Name_Digit1: input_Digit1,
		Name_LMB:    input_LMB,
		Name_RMB:    input_RMB,
	}
)

func Test_NewInputManager(t *testing.T) {
	inputManager := NewInputManager(&testMapping)

	expectedKeyState := map[string]KeyState{
		Name_KeyA:   KeyState_Up,
		Name_Digit1: KeyState_Up,
		Name_LMB:    KeyState_Up,
		Name_RMB:    KeyState_Up,
	}

	assert.Equal(t, expectedTransformedMapping, inputManager.mapping)
	assert.Equal(t, expectedKeyState, inputManager.inputState)
}

func Test_GetInputState(t *testing.T) {
	inputManager := NewInputManager(&testMapping)
	t.Run("invalid input key", func(t *testing.T) {
		state, err := inputManager.GetInputState("NonExistentKey")
		assert.Equal(t, ErrKeyNotFound, err)
		assert.Equal(t, KeyState_Null, state)
	})
	t.Run("valid input key", func(t *testing.T) {
		state, err := inputManager.GetInputState(Name_KeyA)
		assert.Nil(t, err)
		assert.Equal(t, KeyState_Up, state)
	})
}

func Test_setInputState(t *testing.T) {
	inputManager := NewInputManager(&testMapping)
	t.Run("set null state", func(t *testing.T) {
		err := inputManager.SetInputState(Name_KeyA, KeyState_Null)
		assert.Equal(t, ErrKeyStateNull, err)
		state, err := inputManager.GetInputState(Name_KeyA)
		assert.Nil(t, err)
		assert.Equal(t, KeyState_Up, state)
	})
	t.Run("invalid input key", func(t *testing.T) {
		key := "NonExistentKey"
		err := inputManager.SetInputState(key, KeyState_Up)
		assert.Equal(t, ErrKeyNotFound, err)
		state, err := inputManager.GetInputState(key)
		assert.Equal(t, ErrKeyNotFound, err)
		assert.Equal(t, KeyState_Null, state)
	})
	t.Run("valid input key", func(t *testing.T) {
		err := inputManager.SetInputState(Name_KeyA, KeyState_Down)
		assert.Nil(t, err)
		state, err := inputManager.GetInputState(Name_KeyA)
		assert.Nil(t, err)
		assert.Equal(t, KeyState_Down, state)
	})
}

func Test_GetAllStates(t *testing.T) {
	inputManager := NewInputManager(&testMapping)

	state := inputManager.GetAllStates()

	expectedKeyState := map[string]KeyState{
		Name_KeyA:   KeyState_Up,
		Name_Digit1: KeyState_Up,
		Name_LMB:    KeyState_Up,
		Name_RMB:    KeyState_Up,
	}

	assert.Equal(t, expectedKeyState, state)
}
