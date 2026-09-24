package component

import (
	"encoding/json"
	"fmt"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
)

var ErrColorValueRangeOutOfBounds = fmt.Errorf("ComponentError: Color values should be in the range of %d to %d", base.ColorValueRangeMin, base.ColorValueRangeMax)

const (
	// default value of each shade for the base color object
	defaultShadeValue = 0

	// default value of the alpha attribute
	defaultAlphaValue = 255
)

// Color is used to define the color of the component.
// It is based off the RGBA attributes
type Color struct {
	r int // red shade
	g int // gree shade
	b int // blue shade
	a int // alpha -> used for transparency
}

// intermediary structure of Color used for marshalling and unmarshalling component details
type colorJSON struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	A int `json:"a"`
}

// internal function used to register the base component copy with the componentRegistry
func newBaseColor() *Color {
	return &Color{
		r: defaultShadeValue,
		g: defaultShadeValue,
		b: defaultShadeValue,
		a: defaultAlphaValue,
	}
}

func NewColor(r int, g int, b int, a int) *Color {
	color := newBaseColor()
	color.Set(r, g, b, a)

	return color
}

// returns the shade and transparency of the color
//
// return format (r, g, b, a)
func (c *Color) Get() (int, int, int, int) {
	return c.r, c.g, c.b, c.a
}

// set the shade and transparency
func (c *Color) Set(r int, g int, b int, a int) {
	c.r = min(base.ColorValueRangeMax, max(r, base.ColorValueRangeMin))
	c.g = min(base.ColorValueRangeMax, max(g, base.ColorValueRangeMin))
	c.b = min(base.ColorValueRangeMax, max(b, base.ColorValueRangeMin))
	c.a = min(base.ColorValueRangeMax, max(a, base.ColorValueRangeMin))
}

// Returns the name of the component
func (c *Color) GetComponentName() string {
	return base.ComponentName_Color
}

// Returns a snapshot of the complete data stored in the component
func (c *Color) GetComponentDetails() ComponentDetails {
	data := colorJSON{
		R: c.r,
		G: c.g,
		B: c.b,
		A: c.a,
	}
	byteData, _ := json.Marshal(data)
	return ComponentDetails{
		Name: c.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

func validateRange(val int) error {
	if val > base.ColorValueRangeMax || val < base.ColorValueRangeMin {
		return ErrColorValueRangeOutOfBounds
	}
	return nil
}

// Build component from provided details
func (c *Color) BuildFromDetails(data json.RawMessage) error {
	var componentData colorJSON
	err := json.Unmarshal(data, &componentData)
	if err != nil {
		return err
	}

	if err := validateRange(componentData.R); err != nil {
		return err
	}
	if err := validateRange(componentData.G); err != nil {
		return err
	}
	if err := validateRange(componentData.B); err != nil {
		return err
	}
	if err := validateRange(componentData.A); err != nil {
		return err
	}

	c.r = componentData.R
	c.g = componentData.G
	c.b = componentData.B
	c.a = componentData.A

	return nil
}

var _ Component = &Color{}
