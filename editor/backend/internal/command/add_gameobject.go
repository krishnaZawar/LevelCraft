package command

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// AddGameobjectCommand requests the creation of a new gameobject in the scene
type AddGameobjectCommand struct{}

// Returns the name of the command
func (agc *AddGameobjectCommand) GetCommandName() string {
	return base.Command_AddGameobject
}

// Breaks the command down into events for further processing
func (agc *AddGameobjectCommand) Handle() []models.Event {
	return []models.Event{event.NewAddGameobjectEvent()}
}

// Factory class to convert incoming data into AddGameobjectCommand type
type AddGameobjectCommandFactory struct{}

func NewAddGameobjectCommandFactory() *AddGameobjectCommandFactory {
	return &AddGameobjectCommandFactory{}
}

// Handles the conversion of raw bytes into meaningful AddGameobjectCommand
func (agcf *AddGameobjectCommandFactory) NewCommand(details json.RawMessage) (models.Command, error) {
	return &AddGameobjectCommand{}, nil
}

var _ models.Command = &AddGameobjectCommand{}
var _ models.CommandFactory = &AddGameobjectCommandFactory{}
