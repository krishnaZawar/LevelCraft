package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/controlmanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
)

// SetExecutionStateCommand updates and manages the control managers internal state to control frame execution
type SetExecutionStateCommand struct {
	controlManager *controlmanager.ControlManager
	ExecutionState controlmanager.ExecutionState `json:"executionState"`
}

// returns the command name of SetExecutionStateCommand
func (sesc *SetExecutionStateCommand) GetCommandName() string {
	return base.Command_SetExecutionState
}

// handle sets the execution state for the control manager
func (sesc *SetExecutionStateCommand) Handle() error {
	err := sesc.controlManager.SetExecutionState(sesc.ExecutionState)
	if err != nil {
		ls.ErrorWith(err).Msgf("failed to set execution state to %s", sesc.ExecutionState)
		return err
	}
	ls.Info().Msgf("execution state set to %s", sesc.ExecutionState)
	return nil
}

// CommandFactory to convert request to SetExecutionStateCommand type
type SetExecutionStateCommandFactory struct {
	controlManager *controlmanager.ControlManager // used to inject the controlmanager into the command
}

func NewSetExecutionStateCommandFactory(controlManager *controlmanager.ControlManager) *SetExecutionStateCommandFactory {
	return &SetExecutionStateCommandFactory{
		controlManager: controlManager,
	}
}

// converts the incoming commandRequest to SetExecutionStateCommand
func (sescf *SetExecutionStateCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *SetExecutionStateCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	comm.controlManager = sescf.controlManager
	return comm, nil
}

var _ models.Command = &SetExecutionStateCommand{}
var _ models.CommandFactory = &SetExecutionStateCommandFactory{}
