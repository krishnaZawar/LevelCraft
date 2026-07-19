package component

import "github.com/krishnaZawar/LevelCraft/utils/component/base"

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
func (t *Transform) GetComponentDetails() map[string]interface{} {
	return map[string]interface{}{
		"x": t.x,
		"y": t.y,
		"w": t.w,
		"h": t.h,
	}
}

// Build component from provided details
func (t *Transform) BuildFromDetails(data map[string]interface{}) {
	if v, ok := data["x"]; ok {
		switch n := v.(type) {
		case int:
			t.x = n
		case float64:
			t.x = int(n)
		}
	}

	if v, ok := data["y"]; ok {
		switch n := v.(type) {
		case int:
			t.y = n
		case float64:
			t.y = int(n)
		}
	}

	if v, ok := data["w"]; ok {
		switch n := v.(type) {
		case int:
			t.w = n
		case float64:
			t.w = int(n)
		}
	}

	if v, ok := data["h"]; ok {
		switch n := v.(type) {
		case int:
			t.h = n
		case float64:
			t.h = int(n)
		}
	}
}

var _ Component = &Transform{}
