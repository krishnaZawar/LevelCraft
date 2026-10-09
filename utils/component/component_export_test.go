package component

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/basecomp"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
	"github.com/stretchr/testify/assert"
)

func contains(arr []string, value string) bool {
	for _, val := range arr {
		if val == value {
			return true
		}
	}
	return false
}

func Test_GetComponentName(t *testing.T) {
	tests := []struct {
		comp         models.Component
		expectedName string
	}{
		{
			comp:         basecomp.NewBaseTransform(),
			expectedName: base.ComponentName_Transform,
		},
		{
			comp:         basecomp.NewBaseColor(),
			expectedName: base.ComponentName_Color,
		},
	}

	for _, tt := range tests {
		name := tt.comp.GetComponentName()
		assert.Equal(t, tt.expectedName, name)
	}
}

func Test_ComponentList(t *testing.T) {
	listLen := 2
	assert.Equal(t, listLen, len(ComponentList))

	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Transform))
	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Color))
}
