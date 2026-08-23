package inputmanager

import (
	"errors"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/logger"
	"github.com/krishnaZawar/LevelCraft/utils/input"
)

var ls = logger.Get()

var (
	ErrKeyNotFound  = errors.New("error: Key does not exist in the mapping")
	ErrKeyStateNull = errors.New("error: Key state null cannot be set on a key")
)

// Enum to define state transitions for input
type KeyState int

const (
	KeyState_Up KeyState = iota
	KeyState_Down
	KeyState_Null
)

// manages the input interactions and state and storage
//
// all the interactions pass through the inputManager
type InputManager struct {
	mapping    map[string]input.InputData // mapping of names -> codes
	inputState map[string]KeyState        // store managing input state for all keys
}

// map input names -> input data
func mapInputs(inputs []input.InputData, mapping map[string]input.InputData) {
	for _, input := range inputs {
		mapping[input.Name] = input
	}
}

// map the input names -> input data
func transformMapping(inputMapping *input.InputMapping) map[string]input.InputData {
	mapping := map[string]input.InputData{}
	mapInputs(inputMapping.Keyboard, mapping)
	mapInputs(inputMapping.Mouse, mapping)
	return mapping
}

// inits key states with default states
func initKeyStates(mapping map[string]input.InputData) map[string]KeyState {
	inputState := map[string]KeyState{}
	for _, input := range mapping {
		inputState[input.Name] = KeyState_Up
	}
	return inputState
}

// creates a new inputManager with the loaded mapping and default states
func NewInputManager(mapping *input.InputMapping) *InputManager {
	transformedMapping := transformMapping(mapping)
	return &InputManager{
		mapping:    transformedMapping,
		inputState: initKeyStates(transformedMapping),
	}
}

// get the input state for the key with name keyName
func (im *InputManager) GetInputState(keyName string) (KeyState, error) {
	data, found := im.mapping[keyName]
	if !found {
		return KeyState_Null, ErrKeyNotFound
	}
	return im.inputState[data.Name], nil
}

// set input state for the key with name keyName
func (im *InputManager) SetInputState(keyName string, state KeyState) error {
	if state == KeyState_Null {
		return ErrKeyStateNull
	}
	data, found := im.mapping[keyName]
	if !found {
		return ErrKeyNotFound
	}
	im.inputState[data.Name] = state
	return nil
}

// get all the input states
func (im *InputManager) GetAllStates() map[string]KeyState {
	return im.inputState
}

// creates a new input manager object with the actual input mapping data
func newSingleton() *InputManager {
	mapping, err := input.LoadMapping()
	if err != nil {
		ls.ErrorWith(err).Msg("failed to load input mapping")
		return &InputManager{}
	}
	return NewInputManager(mapping)
}

var inputManager *InputManager = newSingleton()

func Get() *InputManager {
	return inputManager
}
