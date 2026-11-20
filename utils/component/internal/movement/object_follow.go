package movement

import (
	"encoding/json"
	"fmt"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
)

var ErrNegativeStopDistance = fmt.Errorf("ComponentError: stopDistance should be not be a negative value")

// Lets an object follow or move toward another object
// ObjectFollow works from the center of the follower object to the center of the target object
type ObjectFollow struct {
	speed float64 // speed with which the object should move

	objectID string // ID of the object that it will follow

	// Limits how close the follower approaches the target.
	// It doesn't maintain a fixed gap, so the target can still move
	// closer while the follower holds its position.
	// Prevents overshooting and jittering around the target.
	stopDistance float64
}

// intermediary structure of ObjectFollow used for marshalling and unmarshalling component details
type objectFollowJSON struct {
	Speed        float64 `json:"speed"`
	ObjectID     string  `json:"objectID,omitempty"`
	StopDistance float64 `json:"stopDistance,omitempty"`
}

// internal function used to register the base component copy with the componentRegistry
func NewBaseObjectFollow() *ObjectFollow {
	return &ObjectFollow{
		speed: defaultSpeed,
	}
}

func NewObjectFollow(speed float64, objectID string, stopDist float64) *ObjectFollow {
	return &ObjectFollow{
		speed:        speed,
		objectID:     objectID,
		stopDistance: stopDist,
	}
}

// Sets the stop distance of the follow object
func (of *ObjectFollow) SetStopDistance(dist float64) {
	of.stopDistance = max(0, dist)
}

// returns the target objects id
func (of *ObjectFollow) GetObjectID() string {
	return of.objectID
}

// returns the speed with which the object will follow the target
func (of *ObjectFollow) GetSpeed() float64 {
	return of.speed
}

// returns the stop distance of the follow object
func (of *ObjectFollow) GetStopDistance() float64 {
	return of.stopDistance
}

// Returns the name of the component
func (of *ObjectFollow) GetComponentName() string {
	return base.ComponentName_ObjectFollow
}

// Returns a snapshot of the complete data stored in the component
func (of *ObjectFollow) GetComponentDetails() models.ComponentDetails {
	data := objectFollowJSON{
		Speed:        of.speed,
		ObjectID:     of.objectID,
		StopDistance: of.stopDistance,
	}
	byteData, _ := json.Marshal(data)
	return models.ComponentDetails{
		Name: of.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

// Build component from provided details
func (of *ObjectFollow) BuildFromDetails(data json.RawMessage) error {
	var componentData objectFollowJSON
	err := json.Unmarshal(data, &componentData)
	if err != nil {
		return err
	}

	if componentData.StopDistance < 0 {
		return ErrNegativeStopDistance
	}

	of.speed = componentData.Speed
	of.objectID = componentData.ObjectID
	of.stopDistance = componentData.StopDistance

	return nil
}

// does not handle any input
func (of *ObjectFollow) HandlesInput(input string) bool {
	return false
}

var _ models.Component = &ObjectFollow{}
