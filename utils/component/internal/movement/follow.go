package movement

import (
	"encoding/json"
	"math"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
)

const (
	// default speed value the component holds on init
	defaultFollowSpeed = 100
)

// What a Follow component is chasing.
const (
	FollowTarget_Object = "object" // another gameobject, by id
	FollowTarget_Point  = "point"  // a fixed position in the world
)

// Follow moves the object towards something else, with no input.
//
// Chasing a gameobject and holding a position are the same motion with a
// different target, so they are one component rather than two.
//
// A Follow with no target is inert rather than invalid: the component has to
// be addable in the editor before a target has been picked, and defaulting to
// the origin would drag every new object to the corner of the scene.
type Follow struct {
	speed      float64 // the speed with which the object should move
	targetMode string  // whether the target is an object or a fixed point

	targetID string // id of the gameobject being followed, in object mode

	targetX float64 // x of the position being held, in point mode
	targetY float64 // y of the position being held, in point mode

	// how close this object will bring itself to the target. It bounds the
	// follower's own approach, not the gap: a target that moves can still
	// walk straight into the follower, which holds its ground rather than
	// backing away. Without it the follower overshoots the target every frame
	// and turns back, which reads as a jitter on top of it.
	stopDistance float64
}

// intermediary structure of Follow used for marshalling and unmarshalling component details
type followJSON struct {
	Speed        float64 `json:"speed"`
	TargetMode   string  `json:"targetMode"`
	TargetID     string  `json:"targetId,omitempty"`
	TargetX      float64 `json:"targetX"`
	TargetY      float64 `json:"targetY"`
	StopDistance float64 `json:"stopDistance"`
}

// internal function used to register the base component copy with the componentRegistry
func NewBaseFollow() *Follow {
	return &Follow{
		speed:      defaultFollowSpeed,
		targetMode: FollowTarget_Object,
	}
}

func NewFollow(speed float64, targetMode string, stopDistance float64) *Follow {
	follow := NewBaseFollow()
	follow.speed = speed
	follow.targetMode = targetMode
	follow.stopDistance = stopDistance
	return follow
}

// get the configured speed for the component
func (f *Follow) GetSpeed() float64 {
	return f.speed
}

// whether this component is chasing a gameobject or holding a position
func (f *Follow) GetTargetMode() string {
	return f.targetMode
}

// id of the gameobject being followed, empty unless in object mode
func (f *Follow) GetTargetID() string {
	if f.targetMode != FollowTarget_Object {
		return ""
	}
	return f.targetID
}

// the position being held, and whether this component is holding one at all
func (f *Follow) GetTargetPoint() (float64, float64, bool) {
	if f.targetMode != FollowTarget_Point {
		return 0, 0, false
	}
	return f.targetX, f.targetY, true
}

// how close to the target counts as arrived
func (f *Follow) GetStopDistance() float64 {
	return f.stopDistance
}

// reports whether this component has been given something to follow
func (f *Follow) HasTarget() bool {
	if f.targetMode == FollowTarget_Point {
		return true
	}
	return f.targetID != ""
}

// Returns the velocity that carries the object towards the target, and zero
// once it is within the stop distance.
//
// dt lets the last step land exactly on the stop distance rather than past it.
func (f *Follow) VelocityTowards(x, y, targetX, targetY, dt float64) (float64, float64) {
	if dt <= 0 {
		return 0, 0
	}

	dx, dy := targetX-x, targetY-y
	distance := math.Hypot(dx, dy)

	remaining := distance - f.stopDistance
	if remaining <= 0 || distance == 0 {
		return 0, 0
	}

	// Close enough that a full step would carry the object past where it
	// should stop, so travel exactly what is left.
	if f.speed*dt >= remaining {
		return (dx / distance) * remaining / dt, (dy / distance) * remaining / dt
	}

	return (dx / distance) * f.speed, (dy / distance) * f.speed
}

// Returns the name of the component
func (f *Follow) GetComponentName() string {
	return base.ComponentName_Follow
}

// Returns a snapshot of the complete data stored in the component
func (f *Follow) GetComponentDetails() models.ComponentDetails {
	data := followJSON{
		Speed:        f.speed,
		TargetMode:   f.targetMode,
		TargetID:     f.targetID,
		TargetX:      f.targetX,
		TargetY:      f.targetY,
		StopDistance: f.stopDistance,
	}
	byteData, _ := json.Marshal(data)
	return models.ComponentDetails{
		Name: f.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

// Build component from provided details
func (f *Follow) BuildFromDetails(data json.RawMessage) error {
	var componentData followJSON
	if err := json.Unmarshal(data, &componentData); err != nil {
		return err
	}

	if componentData.Speed < 0 {
		return base.ErrNegativeSpeed
	}
	if componentData.StopDistance < 0 {
		return base.ErrNegativeStopDistance
	}

	switch componentData.TargetMode {
	case FollowTarget_Object, FollowTarget_Point:
	case "":
		// Written before the mode existed, or by a client that left it out.
		componentData.TargetMode = FollowTarget_Object
	default:
		return base.NewErrUnknownFollowTarget(componentData.TargetMode)
	}

	// The id is meaningless in point mode, so it is dropped rather than kept
	// around to be followed again if the mode is switched back.
	if componentData.TargetMode != FollowTarget_Object {
		componentData.TargetID = ""
	}

	f.speed = componentData.Speed
	f.targetMode = componentData.TargetMode
	f.targetID = componentData.TargetID
	f.targetX = componentData.TargetX
	f.targetY = componentData.TargetY
	f.stopDistance = componentData.StopDistance

	return nil
}

// the target is fixed, so no input affects it
func (f *Follow) HandlesInput(_ string) bool {
	return false
}

var _ models.Component = &Follow{}
