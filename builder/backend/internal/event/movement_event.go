package event

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// Every movement related component just emits the MoveEvent for handling movements for the object
type MoveEvent struct {
	GameobjectID string // objectID of the object to be moved

	MoveX float64 // movement for x direction
	MoveY float64 // movement for y direction
}

func (me *MoveEvent) GetEventName() string {
	return base.Event_MoveEvent
}

var _ models.Event = &MoveEvent{}
