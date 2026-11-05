package event

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// Indicates that some kind of collision occurred between the given pair of objects
type CollisionOccuredEvent struct {
	Obj1_ID string // objID of the first object in the pair
	Obj2_ID string // objID of the second object in the pair
}

func (coe *CollisionOccuredEvent) GetEventName() string {
	return base.Event_CollisionOccuredEvent
}

var _ models.Event = &CollisionOccuredEvent{}
