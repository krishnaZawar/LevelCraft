package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// DeleteGameobjectCommand requests the deletion of the gameobject with "id" in the game scene
type DeleteGameobjectCommand struct {
	ID string `json:"id"` // the unique identifier used for deletion
}

// Returns the name of the command
func (dgc *DeleteGameobjectCommand) GetCommandName() string {
	return base.Command_DeleteGameobject
}

// Breaks the command down into events for further processing
func (dgc *DeleteGameobjectCommand) Handle() []models.Event {
	return []models.Event{event.NewDeleteGameobjectEvent(dgc.ID)}
}

// Factory class to convert incoming data into DeleteGameobjectCommand type
type DeleteGameobjectCommandFactory struct{}

func NewDeleteGameobjectCommandFactory() *DeleteGameobjectCommandFactory {
	return &DeleteGameobjectCommandFactory{}
}

// Handles the conversion of raw bytes into meaningful DeleteGameobjectCommand
func (dgcf *DeleteGameobjectCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	var comm *DeleteGameobjectCommand
	if err := json.Unmarshal(details, &comm); err != nil {
		return nil, err
	}
	return comm, nil
}

var _ models.Command = &DeleteGameobjectCommand{}
var _ models.CommandFactory = &DeleteGameobjectCommandFactory{}
