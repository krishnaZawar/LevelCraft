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
	c.r = r
	c.g = g
	c.b = b
	c.a = a
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
func (c *Color) BuildFromDetails(data map[string]interface{}) {
	if v, ok := data["r"]; ok {
		switch n := v.(type) {
		case int:
			c.r = n
		case float64:
			c.r = int(n)
		}
	}

	if v, ok := data["g"]; ok {
		switch n := v.(type) {
		case int:
			c.g = n
		case float64:
			c.g = int(n)
		}
	}

	if v, ok := data["b"]; ok {
		switch n := v.(type) {
		case int:
			c.b = n
		case float64:
			c.b = int(n)
		}
	}

	if v, ok := data["a"]; ok {
		switch n := v.(type) {
		case int:
			c.a = n
		case float64:
			c.a = int(n)
		}
	}
}

var _ Component = &Color{}
