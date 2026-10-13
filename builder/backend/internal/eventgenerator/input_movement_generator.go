package eventgenerator

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// Generates the Events that the InputMovement component should emit
// if any key that interacts with it is pressed
type InputMovementEventGenerator struct {
	gsm *gamestatemanager.GameStateManager
}

func NewInputMovementEventGenerator(gsm *gamestatemanager.GameStateManager) *InputMovementEventGenerator {
	return &InputMovementEventGenerator{
		gsm: gsm,
	}
}

// Generates the MoveEvent event for the object to update its Transform and draw the new gamestatee
func (imeg *InputMovementEventGenerator) GenerateEvents(ctx EventGeneratorCtx) ([]models.Event, error) {
	if ctx.Component == nil {
		return nil, base.ErrNilComponentPassed
	}
	inputComp, ok := ctx.Component.(*component.InputMovement)
	if !ok {
		return nil, base.NewErrComponentTypeMisMatch(compbase.ComponentName_InputMovement, ctx.Component.GetComponentName())
	}

	_, ok = imeg.gsm.GetGameobject(ctx.GameobjectID)
	if !ok {
		return nil, base.NewErrGameobjectNotFound(ctx.GameobjectID)
	}

	xdir, ydir := inputComp.GetDirections(ctx.InputKey)

	return []models.Event{
		&event.MoveEvent{
			GameobjectID: ctx.GameobjectID,
			MoveX:        float64(xdir) * inputComp.GetSpeed(),
			MoveY:        float64(ydir) * inputComp.GetSpeed(),
		},
	}, nil
}

var _ EventGenerator = &InputMovementEventGenerator{}
