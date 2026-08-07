package component

import "github.com/krishnaZawar/LevelCraft/utils/component/base"

const (
	// default value of each shade for the base color object
	defaultShadeValue = 0

	//default value of the alpha attribute
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
func (c *Color) GetComponentDetails() map[string]interface{} {
	return map[string]interface{}{
		"r": c.r,
		"g": c.g,
		"b": c.b,
		"a": c.a,
	}
}

// Build component from provided details
func (c *Color) BuildFromDetails(data map[string]interface{}) error {
	temp := *c
	if v, ok := data["r"]; ok {
		switch n := v.(type) {
		case int:
			temp.r = n
		case float64:
			temp.r = int(n)
		default:
			return base.ErrExpectedInteger
		}
		if temp.r < base.ColorValueRangeMin || temp.r > base.ColorValueRangeMax {
			return base.ErrColorValueRangeOutOfBounds
		}
	}

	if v, ok := data["g"]; ok {
		switch n := v.(type) {
		case int:
			temp.g = n
		case float64:
			temp.g = int(n)
		default:
			return base.ErrExpectedInteger
		}
		if temp.g < base.ColorValueRangeMin || temp.g > base.ColorValueRangeMax {
			return base.ErrColorValueRangeOutOfBounds
		}
	}

	if v, ok := data["b"]; ok {
		switch n := v.(type) {
		case int:
			temp.b = n
		case float64:
			temp.b = int(n)
		default:
			return base.ErrExpectedInteger
		}
		if temp.b < base.ColorValueRangeMin || temp.b > base.ColorValueRangeMax {
			return base.ErrColorValueRangeOutOfBounds
		}
	}

	if v, ok := data["a"]; ok {
		switch n := v.(type) {
		case int:
			temp.a = n
		case float64:
			temp.a = int(n)
		default:
			return base.ErrExpectedInteger
		}
		if temp.a < base.ColorValueRangeMin || temp.a > base.ColorValueRangeMax {
			return base.ErrColorValueRangeOutOfBounds
		}
	}

	*c = temp

	return nil
}

var _ Component = &Color{}
