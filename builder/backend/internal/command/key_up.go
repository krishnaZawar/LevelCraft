package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// KeyUpCommand indicates that a keyboard key was released and inputManager needs to be updated
type KeyUpCommand struct {
	// injected managers for handling and breaking down of commands into events
	inputManager *inputmanager.InputManager

	KeyName string `json:"keyName"` // identifier of the key released
}

// returns the command name of KeyUpCommand
func (kuc *KeyUpCommand) GetCommandName() string {
	return base.Command_KeyUp
}

// break down the command into events based on the game semantics
func (kuc *KeyUpCommand) Handle() []models.Event {
	// pass through the interpreter and handle command
	err := kuc.inputManager.SetInputState(kuc.KeyName, inputmanager.KeyState_Up)
	if err != nil {
		ls.ErrorWith(err).Msgf("failed to set key state up for %s", kuc.KeyName)
	}
	return []models.Event{}
}

// CommandFactory to convert request to KeyUpCommand type
type KeyUpCommandFactory struct {
	inputmanager *inputmanager.InputManager // used to inject the inputmanager into the command
}

func NewKeyUpCommandFactory(inputManager *inputmanager.InputManager) *KeyUpCommandFactory {
	return &KeyUpCommandFactory{
		inputmanager: inputManager,
	}
}

// converts the incoming commandRequest to KeyUpCommand
func (kucf *KeyUpCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *KeyUpCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	comm.inputManager = kucf.inputmanager
	return comm, nil
}

var _ models.Command = &KeyUpCommand{}
var _ models.CommandFactory = &KeyUpCommandFactory{}
