package basecomp

import (
	"encoding/json"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
)

const (
	// default value of each attribute for the base transform object
	defaultTransformValue = 100.0
)

// Transform is used to determine the position and dimension of any object in the game scene
type Transform struct {
	x float64 // x coordinate of the object
	y float64 // y coordinate of the object
	w float64 // width of the object
	h float64 // height of the object
}

// intermediary structure of Transform used for marshalling and unmarshalling component details
type transformJSON struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

// internal function used to register the base component copy with the componentRegistry
func NewBaseTransform() *Transform {
	return &Transform{
		x: defaultTransformValue,
		y: defaultTransformValue,
		w: defaultTransformValue,
		h: defaultTransformValue,
	}
}

func NewTransform(x float64, y float64, w float64, h float64) *Transform {
	transform := NewBaseTransform()
	transform.UpdatePosition(x, y)
	transform.UpdateDimension(w, h)
	return transform
}

// Returns the position of the object
//
// Return type is (float, float) indicating (x coordinate, y coordinate)
func (t *Transform) GetPosition() (float64, float64) {
	return t.x, t.y
}

// Returns the dimension of the object
//
// Return type is (float, float) indicating (width, height)
func (t *Transform) GetDimension() (float64, float64) {
	return t.w, t.h
}

// Updates the position of the object
func (t *Transform) UpdatePosition(x float64, y float64) {
	t.x = x
	t.y = y
}

// Updates the dimension of the object
func (t *Transform) UpdateDimension(w float64, h float64) {
	t.w = w
	t.h = h
}

// Returns the name of the component
func (t *Transform) GetComponentName() string {
	return base.ComponentName_Transform
}

// Returns a snapshot of the complete data stored in the component
func (t *Transform) GetComponentDetails() models.ComponentDetails {
	data := transformJSON{
		X: t.x,
		Y: t.y,
		W: t.w,
		H: t.h,
	}
	byteData, _ := json.Marshal(data)
	return models.ComponentDetails{
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

// no interaction possible
func (t *Transform) HandlesInput(input string) bool {
	return false
}

var _ models.Component = &Transform{}
