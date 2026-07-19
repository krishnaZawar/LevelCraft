package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// UpdateComponentCommand requests the updation of the component, if exists, for a specific gameobject in the scene
// update happens only if the component exists
type UpdateComponentCommand struct {
	ID            string                 `json:"id"`            // unique identifier of the gameobject from which the component is updated
	ComponentName string                 `json:"componentName"` // name of the component to be updated
	Data          map[string]interface{} `json:"data"`          // data of the new command
}

// Returns the name of the command
func (ucc *UpdateComponentCommand) GetCommandName() string {
	return base.Command_UpdateComponent
}

// Breaks the command down into events for further processing
func (ucc *UpdateComponentCommand) Handle() []models.Event {
	return []models.Event{event.NewUpdateComponentEvent(ucc.ID, ucc.ComponentName, ucc.Data)}
}

// Factory class to convert incoming data into UpdateComponentCommand type
type UpdateComponentCommandFactory struct{}

func NewUpdateComponentCommandFactory() *UpdateComponentCommandFactory {
	return &UpdateComponentCommandFactory{}
}

// Handles the conversion of raw bytes into meaningful UpdateComponentCommand
func (uccf *UpdateComponentCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *UpdateComponentCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	return comm, nil
}

var _ models.Command = &UpdateComponentCommand{}
var _ models.CommandFactory = &UpdateComponentCommandFactory{}
