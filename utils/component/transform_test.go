package component

import (
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

func Test_GetTransformComponentName(t *testing.T) {
	tr := NewTransform(50, 50, 50, 50)

	assert.Equal(t, base.ComponentName_Transform, tr.GetComponentName())
}

func Test_GetTransformComponentDetails(t *testing.T) {
	x, y, w, h := 50, 50, 50, 50
	tr := NewTransform(x, y, w, h)

	data := tr.GetComponentDetails()

	assert.Equal(t, x, data["x"].(int))
	assert.Equal(t, y, data["y"].(int))
	assert.Equal(t, w, data["w"].(int))
	assert.Equal(t, h, data["h"].(int))
}

func Test_BuildFromDetails(t *testing.T) {
	tr := newBaseTransform()

	data := map[string]interface{}{
		"x": 12,
		"w": 10,
	}

	tr.BuildFromDetails(data)

	x, y := tr.GetPosition()
	w, h := tr.GetDimension()
	assert.Equal(t, 12, x)
	assert.Equal(t, defaultTransformValue, y)
	assert.Equal(t, 10, w)
	assert.Equal(t, defaultTransformValue, h)
}
