package component

import (
	"testing"

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
