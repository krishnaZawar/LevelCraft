package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// DetachComponentCommand requests the removal of the component, if exists, from a specific gameobject in the scene
type DetachComponentCommand struct {
	ID            string `json:"id"`            // unique identifier of the gameobject from which the component is removed
	ComponentName string `json:"componentName"` // name of the component to be removed
}

// Returns the name of the command
func (dcc *DetachComponentCommand) GetCommandName() string {
	return base.Command_DetachComponent
}

// Breaks the command down into events for further processing
func (dcc *DetachComponentCommand) Handle() []models.Event {
	return []models.Event{event.NewDetachComponentEvent(dcc.ID, dcc.ComponentName)}
}

// Factory class to convert incoming data into DetachComponentCommand type
type DetachComponentCommandFactory struct{}

func NewDetachComponentCommandFactory() *DetachComponentCommandFactory {
	return &DetachComponentCommandFactory{}
}

// Handles the conversion of raw bytes into meaningful DetachComponentCommand
func (dccf *DetachComponentCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *DetachComponentCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	return comm, nil
}

var _ models.Command = &DetachComponentCommand{}
var _ models.CommandFactory = &DetachComponentCommandFactory{}
