package collisionsystem

import (
	"fmt"

	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/base"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/builder/backend/internal/gamestatemanager"
	"github.com/krishnaZawar/LevelCraft/utils/component"
	compbase "github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/gameobject"
	"github.com/krishnaZawar/LevelCraft/utils/models"
)

// Stores the overlap detected between the 2 objects
// all 4 positive 					-> collision occurred
// atleast 1 zero rest positive		-> edge touching
// atleast 1 negative 				-> no collision
type collisionOverlap struct {
	Left   float64 // left overlap
	Right  float64 // right overlap
	Top    float64 // top overlap
	Bottom float64 // bottom overlap
}

// calculates and returns the collision overlap for the 2 transforms using AABB theorem of geometry
// Here the calculation is done considering the PrimaryObj as the source object for whom overlap is calculated
func calculateOverlap(primaryObj_transform, secondaryObj_transform *component.Transform) collisionOverlap {
	x1, y1 := primaryObj_transform.GetPosition()
	w1, h1 := primaryObj_transform.GetDimension()

	x2, y2 := secondaryObj_transform.GetPosition()
	w2, h2 := secondaryObj_transform.GetDimension()

	return collisionOverlap{
		Left:   x1 + w1 - x2,
		Right:  x2 + w2 - x1,
		Top:    y1 + h1 - y2,
		Bottom: y2 + h2 - y1,
	}
}

// Detects and Resolves collisions between gameobjects
type CollisionSystem struct {
	gsm *gamestatemanager.GameStateManager
}

func NewCollisionSystem(gsm *gamestatemanager.GameStateManager) *CollisionSystem {
	return &CollisionSystem{
		gsm: gsm,
	}
}

// Detects collision between pairs via the AABB theorem of geometry
func (cs *CollisionSystem) DetectCollision(obj1, obj2 *gameobject.Gameobject) bool {
	// return if any object longer exists
	if obj1 == nil || obj2 == nil {
		return false
	}

	// check both have collider and transform component to proceed with detection
	comp, ok := obj1.GetComponent(compbase.ComponentName_Collider)
	if !ok {
		return false
	}
	_, ok = comp.(*component.Collider)
	if !ok {
		return false
	}
	comp, ok = obj1.GetComponent(compbase.ComponentName_Transform)
	if !ok {
		return false
	}
	obj1_transformComp, ok := comp.(*component.Transform)
	if !ok {
		return false
	}

	comp, ok = obj2.GetComponent(compbase.ComponentName_Collider)
	if !ok {
		return false
	}
	_, ok = comp.(*component.Collider)
	if !ok {
		return false
	}
	comp, ok = obj2.GetComponent(compbase.ComponentName_Transform)
	if !ok {
		return false
	}
	obj2_transformComp, ok := comp.(*component.Transform)
	if !ok {
		return false
	}

	// detect if the collision occurred between them or not
	overlap := calculateOverlap(obj1_transformComp, obj2_transformComp)

	xOverlap := min(overlap.Left, overlap.Right)
	yOverlap := min(overlap.Top, overlap.Bottom)

	return xOverlap > 0 && yOverlap > 0
}

// Handle the CollisionOccuredEvent and resolves the collision between them
func (cs *CollisionSystem) Handle(evt models.Event, _ float64) *models.EventResponse {
	collidedEvt, ok := evt.(*event.CollisionOccuredEvent)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgExpectedEventTypeNotFound, base.NewErrEventTypeMismatch(base.Event_CollisionOccuredEvent, evt.GetEventName()))
	}

	// fetch gameobjects
	primary, ok := cs.gsm.GetGameobject(collidedEvt.Obj1_ID)
	if !ok {
		return event.NewNonEmittableResponse(false, base.MsgGameobjectNotFound, base.NewErrGameobjectNotFound(collidedEvt.Obj1_ID))
	}
	secondary, ok := cs.gsm.GetGameobject(collidedEvt.Obj2_ID)
	if !ok {
		return event.NewNonEmittableResponse(false, base.MsgGameobjectNotFound, base.NewErrGameobjectNotFound(collidedEvt.Obj2_ID))
	}

	// fetch Collider and Transform components
	comp, ok := primary.GetComponent(compbase.ComponentName_Collider)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentNotFound(compbase.ComponentName_Collider, collidedEvt.Obj1_ID))
	}
	primary_colliderComp, ok := comp.(*component.Collider)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Collider, comp.GetComponentName()))
	}
	comp, ok = primary.GetComponent(compbase.ComponentName_Transform)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentNotFound(compbase.ComponentName_Transform, collidedEvt.Obj1_ID))
	}
	primary_transformComp, ok := comp.(*component.Transform)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Transform, comp.GetComponentName()))
	}

	comp, ok = secondary.GetComponent(compbase.ComponentName_Collider)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentNotFound(compbase.ComponentName_Collider, collidedEvt.Obj2_ID))
	}
	secondary_colliderComp, ok := comp.(*component.Collider)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Collider, comp.GetComponentName()))
	}
	comp, ok = secondary.GetComponent(compbase.ComponentName_Transform)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentNotFound(compbase.ComponentName_Transform, collidedEvt.Obj2_ID))
	}
	secondary_transformComp, ok := comp.(*component.Transform)
	if !ok {
		return event.NewEmittableResponse(false, base.MsgComponentNotFound, base.NewErrComponentTypeMisMatch(compbase.ComponentName_Transform, comp.GetComponentName()))
	}

	// check if either is movable, then only it makes sense for resolution
	primary_movable := !primary_colliderComp.IsImmovable()
	secondary_movable := !secondary_colliderComp.IsImmovable()

	if !primary_movable && !secondary_movable {
		return event.NewEmittableResponse(true, fmt.Sprintf("collision resolution skipped as objects %s and %s are marked as immovable", collidedEvt.Obj1_ID, collidedEvt.Obj2_ID), nil)
	}

	// calculate collision
	overlap := calculateOverlap(primary_transformComp, secondary_transformComp)

	xResolution := min(overlap.Left, overlap.Right)
	yResolution := min(overlap.Top, overlap.Bottom)

	if xResolution <= 0 || yResolution <= 0 {
		return event.NewEmittableResponse(true, "collision no longer exists, skipping resolution", nil)
	}

	// resolve collision
	xPrimary, yPrimary := primary_transformComp.GetPosition()
	xSecondary, ySecondary := secondary_transformComp.GetPosition()

	var dx, dy float64
	if xResolution < yResolution {
		if xResolution == overlap.Left {
			dx = -overlap.Left
		} else {
			dx = overlap.Right
		}
	} else {
		if yResolution == overlap.Top {
			dy = -overlap.Top
		} else {
			dy = overlap.Bottom
		}
	}

	/*
		Note:
			the resolution is not smoothened by delta time (dt) every frame because:
			- The calculations of resolution are exact and represents the value that is dependent on the frame heavily
				and as it is frame dependent there is no need for dt smoothening here
	*/

	switch {
	case primary_movable && secondary_movable:
		primary_transformComp.UpdatePosition(xPrimary+dx/2, yPrimary+dy/2)
		secondary_transformComp.UpdatePosition(xSecondary-dx/2, ySecondary-dy/2)
	case primary_movable:
		primary_transformComp.UpdatePosition(xPrimary+dx, yPrimary+dy)
	case secondary_movable:
		secondary_transformComp.UpdatePosition(xSecondary-dx, ySecondary-dy)
	}

	return event.NewNonEmittableResponse(true, "collision resolved successfully", nil)
}

var _ models.EventHandler = &CollisionSystem{}

var system *CollisionSystem = NewCollisionSystem(gamestatemanager.Get())

func Get() *CollisionSystem {
	return system
}
