package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/inputmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/logger"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

var ls = logger.Get()

// KeyDownCommand indicates that a keyboard key was pressed and action needs to be taken
type KeyDownCommand struct {
	// injected managers for handling and breaking down of commands into events
	inputManager *inputmanager.InputManager

	KeyName string `json:"keyName"` // identifier of the key pressed
}

// returns the command name of KeyDownCommand
func (kdc *KeyDownCommand) GetCommandName() string {
	return base.Command_KeyDown
}

// break down the command into events based on the game semantics
func (kdc *KeyDownCommand) Handle() []models.Event {
	// pass through the interpreter and handle command
	err := kdc.inputManager.SetInputState(kdc.KeyName, inputmanager.KeyState_Down)
	if err != nil {
		ls.ErrorWith(err).Msgf("failed to set key state down for %s", kdc.KeyName)
	}
	return []models.Event{}
}

// CommandFactory to convert request to KeyDownCommand type
type KeyDownCommandFactory struct {
	inputmanager *inputmanager.InputManager // used to inject the inputmanager into the command
}

func NewKeyDownCommandFactory(inputManager *inputmanager.InputManager) *KeyDownCommandFactory {
	return &KeyDownCommandFactory{
		inputmanager: inputManager,
	}
}

// converts the incoming commandRequest to KeyDownCommand
func (kdcf *KeyDownCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *KeyDownCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	comm.inputManager = kdcf.inputmanager
	return comm, nil
}

var _ models.Command = &KeyDownCommand{}
var _ models.CommandFactory = &KeyDownCommandFactory{}
