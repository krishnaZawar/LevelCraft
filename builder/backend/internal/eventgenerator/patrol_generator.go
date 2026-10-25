package eventgenerator

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// Generates the events the Patrol component emits.
//
// Patrol runs off the frame rather than off input, so this is registered as a
// frame generator: there is no key press to trigger it.
type PatrolEventGenerator struct {
	gsm *gamestatemanager.GameStateManager
}

func NewPatrolEventGenerator(gsm *gamestatemanager.GameStateManager) *PatrolEventGenerator {
	return &PatrolEventGenerator{
		gsm: gsm,
	}
}

// Generates the MoveEvent that carries the object towards its next point
func (peg *PatrolEventGenerator) GenerateEvents(ctx EventGeneratorCtx) ([]models.Event, error) {
	if ctx.Component == nil {
		return nil, base.ErrNilComponentPassed
	}
	patrolComp, ok := ctx.Component.(*component.Patrol)
	if !ok {
		return nil, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Patrol, ctx.Component.GetComponentName())
	}

	obj, ok := peg.gsm.GetGameobject(ctx.GameobjectID)
	if !ok {
		return nil, base.NewErrGameobjectNotFound(ctx.GameobjectID)
	}

	// Where the object is now decides where it goes next, so a Patrol without
	// a Transform has nothing to move.
	transform, ok := transformOf(obj)
	if !ok {
		return nil, base.NewErrComponentNotFound(compbase.ComponentName_Transform, ctx.GameobjectID)
	}

	x, y := transform.GetPosition()
	moveX, moveY := patrolComp.NextVelocity(x, y, ctx.DeltaTime)

	// A route of fewer than two points, or an object already standing on its
	// next point, produces no movement — and no event worth queueing.
	if moveX == 0 && moveY == 0 {
		return nil, nil
	}

	return []models.Event{
		&event.MoveEvent{
			GameobjectID: ctx.GameobjectID,
			MoveX:        moveX,
			MoveY:        moveY,
		},
	}, nil
}

var _ EventGenerator = &PatrolEventGenerator{}
