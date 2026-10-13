package eventhandler

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// handles the move event and updates the game state accordingly
// This is handled by the Transform Component as that is responsible for the object positioning
type MoveEventHandler struct {
	gsm *gamestatemanager.GameStateManager
}

func NewMoveEventHandler(gsm *gamestatemanager.GameStateManager) *MoveEventHandler {
	return &MoveEventHandler{
		gsm: gsm,
	}
}

// Handles the MoveEvent and Updates the transform component for that object
func (meh *MoveEventHandler) Handle(evt models.Event, dt float64) *models.EventResponse {
	moveEvt, ok := evt.(*event.MoveEvent)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgExpectedEventTypeNotFound, base.NewErrEventTypeMismatch(base.Event_MoveEvent, evt.GetEventName()))
	}

	obj, ok := meh.gsm.GetGameobject(moveEvt.GameobjectID)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgGameobjectNotFound, base.NewErrGameobjectNotFound(moveEvt.GameobjectID))
	}

	comp, ok := obj.GetComponent(compbase.ComponentName_Transform)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentNotFound(compbase.ComponentName_Transform, moveEvt.GameobjectID))
	}

	transformComp, ok := comp.(*component.Transform)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentNotFound(compbase.ComponentName_Transform, moveEvt.GameobjectID))
	}

	x, y := transformComp.GetPosition()

	newX := x + moveEvt.MoveX*dt
	newY := y + moveEvt.MoveY*dt

	transformComp.UpdatePosition(newX, newY)

	return event.NewEmittableResponse(true, "object position updated successfully", nil)
}

var _ models.EventHandler = &MoveEventHandler{}
