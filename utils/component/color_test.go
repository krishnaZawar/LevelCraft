package component

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_GetColor(t *testing.T) {
	r, g, b, a := 100, 100, 100, 100
	color := NewColor(r, g, b, a)

	actualR, actualG, actualB, actualA := color.Get()

	assert.Equal(t, r, actualR)
	assert.Equal(t, g, actualG)
	assert.Equal(t, b, actualB)
	assert.Equal(t, a, actualA)
}

func Test_SetColor(t *testing.T) {
	color := newBaseColor()

	r, g, b, a := 100, 100, 100, 100

	color.Set(r, g, b, a)

	actualR, actualG, actualB, actualA := color.Get()

	assert.Equal(t, r, actualR)
	assert.Equal(t, g, actualG)
	assert.Equal(t, b, actualB)
	assert.Equal(t, a, actualA)
}

func Test_MarshalAndUnmarshalColor(t *testing.T) {
	r, g, b, a := 100, 120, 140, 160
	expected := NewColor(r, g, b, a)

	found := newBaseColor()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)
	assert.Nil(t, err)

	assert.Equal(t, base.ComponentName_Color, details.Name)
	assert.Equal(t, expected, found)
}

func Test_UnmarshalColor(t *testing.T) {
	t.Run("valid unmarshal", func(t *testing.T) {
		data := json.RawMessage(
			`{
			"r": 100,
			"g": 120,
			"b": 140,
			"a": 160
			}`,
		)
		comp := newBaseColor()
		err := comp.BuildFromDetails(data)
		assert.Nil(t, err)

		assert.Equal(t, 100, comp.r)
		assert.Equal(t, 120, comp.g)
		assert.Equal(t, 140, comp.b)
		assert.Equal(t, 160, comp.a)
	})
	t.Run("invalid unmarshal", func(t *testing.T) {
		data := json.RawMessage(
			`{
			"r": 100,
			"g": 120,
			"b": 140
			"a": 160
			}`,
		)
		comp := newBaseColor()
		err := comp.BuildFromDetails(data)
		assert.NotNil(t, err)

		assert.Equal(t, defaultShadeValue, comp.r)
		assert.Equal(t, defaultShadeValue, comp.g)
		assert.Equal(t, defaultShadeValue, comp.b)
		assert.Equal(t, defaultAlphaValue, comp.a)
	})
	t.Run("r value out of bounds", func(t *testing.T) {
		t.Run("r is below bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": -1,
			"g": 120,
			"b": 140,
			"a": 160
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
		t.Run("r is above bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 256,
			"g": 120,
			"b": 140,
			"a": 160
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
	})

	t.Run("g value out of bounds", func(t *testing.T) {
		t.Run("g is below bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 0,
			"g": -1,
			"b": 140,
			"a": 160
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
		t.Run("g is above bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 255,
			"g": 256,
			"b": 140,
			"a": 160
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
	})

	t.Run("b value out of bounds", func(t *testing.T) {
		t.Run("b is below bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 0,
			"g": 0,
			"b": -1,
			"a": 160
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
		t.Run("b is above bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 255,
			"g": 255,
			"b": 256,
			"a": 160
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
	})

	t.Run("a value out of bounds", func(t *testing.T) {
		t.Run("a is below bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 0,
			"g": 0,
			"b": 0,
			"a": -1
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
		t.Run("a is above bounds", func(t *testing.T) {
			data := json.RawMessage(
				`{
			"r": 255,
			"g": 255,
			"b": 255,
			"a": 256
			}`,
			)
			comp := newBaseColor()
			err := comp.BuildFromDetails(data)
			assert.Equal(t, ErrColorValueRangeOutOfBounds, err)

			assert.Equal(t, defaultShadeValue, comp.r)
			assert.Equal(t, defaultShadeValue, comp.g)
			assert.Equal(t, defaultShadeValue, comp.b)
			assert.Equal(t, defaultAlphaValue, comp.a)
		})
	})
}
