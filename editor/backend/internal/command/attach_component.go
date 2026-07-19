package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// AttachComponentCommand requests the addition of a component to the specified gameobject
type AttachComponentCommand struct {
	ID            string `json:"id"`            // unique identifier of the gameobject to which the component is being added
	ComponentName string `json:"componentName"` // name of the component to be added
}

// Returns the name of the command
func (acc *AttachComponentCommand) GetCommandName() string {
	return base.Command_AttachComponent
}

// Breaks the command down into events for further processing
func (acc *AttachComponentCommand) Handle() []models.Event {
	return []models.Event{event.NewAttachComponentEvent(acc.ID, acc.ComponentName)}
}

// Factory class to convert incoming data into AttachComponentCommand type
type AttachComponentCommandFactory struct{}

func NewAttachComponentCommandFactory() *AttachComponentCommandFactory {
	return &AttachComponentCommandFactory{}
}

// Handles the conversion of raw bytes into meaningful AttachComponentCommand
func (accf *AttachComponentCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *AttachComponentCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	return comm, nil
}

var _ models.Command = &AttachComponentCommand{}
var _ models.CommandFactory = &AttachComponentCommandFactory{}
