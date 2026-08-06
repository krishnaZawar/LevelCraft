package component

import (
	"testing"

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
