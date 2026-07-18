package component

import "github.com/krishnaZawar/LevelCraft/utils/component/base"

// Transform is used to determine the position and dimension of any object in the game scene
type Transform struct {
	x int // x coordinate of the object
	y int // y coordinate of the object
	w int // width of the object
	h int // height of the object
}

func NewTransform(x int, y int, w int, h int) *Transform {
	return &Transform{
		x: x,
		y: y,
		w: w,
		h: h,
	}
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

var _ Component = &Transform{}
