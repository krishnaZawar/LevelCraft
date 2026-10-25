package component

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/basecomp"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/movement"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/stretchr/testify/assert"
)

func Test_NewComponentRegistry(t *testing.T) {
	mapping := &input.InputMapping{}
	compRegistry := NewComponentRegistry(mapping)

	assert.Equal(t, 3, compRegistry.Size())

	var expected models.Component
	var factory models.ComponentFactory
	factory, ok := compRegistry.GetComponent(base.ComponentName_Transform)
	comp := factory()
	expected = basecomp.NewBaseTransform()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)

	factory, ok = compRegistry.GetComponent(base.ComponentName_Color)
	comp = factory()
	expected = basecomp.NewBaseColor()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)

	factory, ok = compRegistry.GetComponent(base.ComponentName_InputMovement)
	comp = factory()
	expected = movement.NewBaseInputMovement(mapping)
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)
}
