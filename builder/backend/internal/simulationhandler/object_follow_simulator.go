package simulationhandler

import (
	"math"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
)

// Runs the simulations the ObjectFollow component is intended to do
type ObjectFollowSimulator struct {
	gsm *gamestatemanager.GameStateManager
}

func NewObjectFollowSimulator(gsm *gamestatemanager.GameStateManager) *ObjectFollowSimulator {
	return &ObjectFollowSimulator{
		gsm: gsm,
	}
}

func (ofs *ObjectFollowSimulator) Simulate(ctx SimulatorCtx) ([]models.Event, error) {
	if ctx.Component == nil {
		return nil, base.ErrNilComponentPassed
	}
	objectFollowComp, ok := ctx.Component.(*component.ObjectFollow)
	if !ok {
		return nil, base.NewErrComponentTypeMisMatch(compbase.ComponentName_ObjectFollow, ctx.Component.GetComponentName())
	}

	followerObj, ok := ofs.gsm.GetGameobject(ctx.GameobjectID)
	if !ok {
		return nil, base.NewErrGameobjectNotFound(ctx.GameobjectID)
	}

	// verify target exists
	targetObj, ok := ofs.gsm.GetGameobject(objectFollowComp.GetObjectID())
	if !ok {
		// return error as nil, because if the targetObj is not found, that is not an error
		// Consider 2 cases:
		// - object died: that is a valid case in game modelling to stop follower simulation
		// - misconfigured object: that is a misconfiguration error that cannot be said as an error
		// Currently will keep it this way, and can update to throw error based on requirements and design
		return []models.Event{}, nil
	}

	// fetch transforms for calculations
	comp, found := followerObj.GetComponent(compbase.ComponentName_Transform)
	if !found {
		return nil, base.NewErrComponentNotFound(compbase.ComponentName_Transform, followerObj.GetID())
	}
	followerTransform, ok := comp.(*component.Transform)
	if !ok {
		return nil, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Transform, comp.GetComponentName())
	}

	comp, found = targetObj.GetComponent(compbase.ComponentName_Transform)
	if !found {
		return nil, base.NewErrComponentNotFound(compbase.ComponentName_Transform, targetObj.GetID())
	}
	targetTransform, ok := comp.(*component.Transform)
	if !ok {
		return nil, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Transform, comp.GetComponentName())
	}

	followerX, followerY := followerTransform.GetPosition()
	followerW, followerH := followerTransform.GetDimension()
	followerCenterX, followerCenterY := followerX+followerW/2, followerY+followerH/2

	targetX, targetY := targetTransform.GetPosition()
	targetW, targetH := targetTransform.GetDimension()
	targetCenterX, targetCenterY := targetX+targetW/2, targetY+targetH/2

	// calculate direction vector
	dx := targetCenterX - followerCenterX
	dy := targetCenterY - followerCenterY

	stopDist := objectFollowComp.GetStopDistance()

	distance := math.Sqrt(dx*dx + dy*dy)
	if distance <= stopDist {
		return []models.Event{}, nil
	}

	speed := objectFollowComp.GetSpeed()
	moveX := 0.0
	moveY := 0.0
	if dx > 0 {
		moveX = 1
	} else if dx < 0 {
		moveX = -1
	}
	if dy > 0 {
		moveY = 1
	} else if dy < 0 {
		moveY = -1
	}

	return []models.Event{
		&event.MoveEvent{
			GameobjectID: ctx.GameobjectID,
			MoveX:        moveX * speed,
			MoveY:        moveY * speed,
		},
	}, nil
}

var _ SimulationHandler = &ObjectFollowSimulator{}
