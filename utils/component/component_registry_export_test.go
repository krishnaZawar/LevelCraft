package component

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/basecomp"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_NewComponentRegistry(t *testing.T) {
	compRegistry := NewComponentRegistry()

	var expected models.Component
	comp, ok := compRegistry.GetComponent(base.ComponentName_Transform)
	expected = basecomp.NewBaseTransform()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)

	comp, ok = compRegistry.GetComponent(base.ComponentName_Color)
	expected = basecomp.NewBaseColor()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)
}
