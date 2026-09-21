package component

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_GetTransform(t *testing.T) {
	x, y, w, h := 100, 100, 100, 100
	tr := NewTransform(x, y, w, h)

	actualX, actualY := tr.GetPosition()
	actualW, actualH := tr.GetDimension()

	assert.Equal(t, x, actualX)
	assert.Equal(t, y, actualY)

	assert.Equal(t, w, actualW)
	assert.Equal(t, h, actualH)
}

func Test_UpdateTransform(t *testing.T) {
	tr := NewTransform(50, 50, 50, 50)

	x, y, w, h := 100, 100, 100, 100

	tr.UpdatePosition(x, y)
	tr.UpdateDimension(w, h)

	actualX, actualY := tr.GetPosition()
	actualW, actualH := tr.GetDimension()

	assert.Equal(t, x, actualX)
	assert.Equal(t, y, actualY)

	assert.Equal(t, w, actualW)
	assert.Equal(t, h, actualH)
}

func Test_MarshalAndUnmarshalTransform(t *testing.T) {
	x, y, w, h := 100, 120, 140, 160
	expected := NewTransform(x, y, w, h)

	found := newBaseTransform()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)
	assert.Nil(t, err)

	assert.Equal(t, base.ComponentName_Transform, details.Name)
	assert.Equal(t, expected, found)
}

func Test_UnmarshalTransform(t *testing.T) {
	t.Run("valid unmarshal", func(t *testing.T) {
		data := json.RawMessage(
			`{
			"x": 100,
			"y": 120,
			"w": 140,
			"h": 160
			}`,
		)
		comp := newBaseTransform()
		err := comp.BuildFromDetails(data)
		assert.Nil(t, err)

		assert.Equal(t, 100, comp.x)
		assert.Equal(t, 120, comp.y)
		assert.Equal(t, 140, comp.w)
		assert.Equal(t, 160, comp.h)
	})
	t.Run("invalid unmarshal", func(t *testing.T) {
		data := json.RawMessage(
			`{
			"x": 100,
			"y": 120,
			"w": 140
			"h": 160
			}`,
		)
		comp := newBaseTransform()
		err := comp.BuildFromDetails(data)
		assert.NotNil(t, err)

		assert.Equal(t, defaultTransformValue, comp.x)
		assert.Equal(t, defaultTransformValue, comp.y)
		assert.Equal(t, defaultTransformValue, comp.w)
		assert.Equal(t, defaultTransformValue, comp.h)
	})
}
