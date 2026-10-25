package eventgenerator

import (
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// the Transform of an object, which every movement generator needs to know
// where it is starting from
func transformOf(obj *gameobject.Gameobject) (*component.Transform, bool) {
	comp, ok := obj.GetComponent(compbase.ComponentName_Transform)
	if !ok {
		return nil, false
	}
	transform, ok := comp.(*component.Transform)
	return transform, ok
}

// Generates the events the Follow component emits.
//
// Like Patrol this runs off the frame rather than off input: a follower chases
// its target whether or not anything is being pressed.
type FollowEventGenerator struct {
	gsm *gamestatemanager.GameStateManager
}

func NewFollowEventGenerator(gsm *gamestatemanager.GameStateManager) *FollowEventGenerator {
	return &FollowEventGenerator{
		gsm: gsm,
	}
}

// resolves where the component is currently heading, which for a followed
// object is wherever that object has moved to this frame
func (feg *FollowEventGenerator) resolveTarget(followComp *component.Follow) (float64, float64, bool) {
	if x, y, isPoint := followComp.GetTargetPoint(); isPoint {
		return x, y, true
	}

	targetID := followComp.GetTargetID()
	if targetID == "" {
		return 0, 0, false
	}

	// A target that has been deleted since the scene was saved leaves the
	// follower standing still rather than failing the frame.
	target, ok := feg.gsm.GetGameobject(targetID)
	if !ok {
		return 0, 0, false
	}
	transform, ok := transformOf(target)
	if !ok {
		return 0, 0, false
	}

	x, y := transform.GetPosition()
	return x, y, true
}

// Generates the MoveEvent that carries the object towards its target
func (feg *FollowEventGenerator) GenerateEvents(ctx EventGeneratorCtx) ([]models.Event, error) {
	if ctx.Component == nil {
		return nil, base.ErrNilComponentPassed
	}
	followComp, ok := ctx.Component.(*component.Follow)
	if !ok {
		return nil, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Follow, ctx.Component.GetComponentName())
	}

	obj, ok := feg.gsm.GetGameobject(ctx.GameobjectID)
	if !ok {
		return nil, base.NewErrGameobjectNotFound(ctx.GameobjectID)
	}
	transform, ok := transformOf(obj)
	if !ok {
		return nil, base.NewErrComponentNotFound(compbase.ComponentName_Transform, ctx.GameobjectID)
	}

	targetX, targetY, hasTarget := feg.resolveTarget(followComp)
	if !hasTarget {
		return nil, nil
	}

	x, y := transform.GetPosition()
	moveX, moveY := followComp.VelocityTowards(x, y, targetX, targetY, ctx.DeltaTime)

	// Already within the stop distance, so there is nothing to emit.
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

var _ EventGenerator = &FollowEventGenerator{}
