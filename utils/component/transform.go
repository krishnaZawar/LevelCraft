package component

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
)

const (
	// default value of each attribute for the base transform object
	defaultTransformValue = 100
)

// Transform is used to determine the position and dimension of any object in the game scene
type Transform struct {
	x int // x coordinate of the object
	y int // y coordinate of the object
	w int // width of the object
	h int // height of the object
}

// intermediary structure of Transform used for marshalling and unmarshalling component details
type transformJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// internal function used to register the base component copy with the componentRegistry
func newBaseTransform() *Transform {
	return &Transform{
		x: defaultTransformValue,
		y: defaultTransformValue,
		w: defaultTransformValue,
		h: defaultTransformValue,
	}
}

func NewTransform(x int, y int, w int, h int) *Transform {
	transform := newBaseTransform()
	transform.UpdatePosition(x, y)
	transform.UpdateDimension(w, h)
	return transform
}

// Returns the position of the object
//
// Return type is (int, int) indicating (x coordinate, y coordinate)
func (t *Transform) GetPosition() (int, int) {
	return t.x, t.y
}

// Returns the dimension of the object
//
// Return type is (int, int) indicating (width, height)
func (t *Transform) GetDimension() (int, int) {
	return t.w, t.h
}

// Updates the position of the object
func (t *Transform) UpdatePosition(x int, y int) {
	t.x = x
	t.y = y
}

// Updates the dimension of the object
func (t *Transform) UpdateDimension(w int, h int) {
	t.w = w
	t.h = h
}

// Returns the name of the component
func (t *Transform) GetComponentName() string {
	return base.ComponentName_Transform
}

// Returns a snapshot of the complete data stored in the component
func (t *Transform) GetComponentDetails() ComponentDetails {
	data := transformJSON{
		X: t.x,
		Y: t.y,
		W: t.w,
		H: t.h,
	}
	byteData, _ := json.Marshal(data)
	return ComponentDetails{
		Name: t.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

// Build component from provided details
func (t *Transform) BuildFromDetails(data json.RawMessage) error {
	var componentData transformJSON
	err := json.Unmarshal(data, &componentData)
	if err != nil {
		return err
	}

	t.x = componentData.X
	t.y = componentData.Y
	t.w = componentData.W
	t.h = componentData.H

	return nil
}

var _ Component = &Transform{}
