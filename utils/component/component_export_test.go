package component

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/basecomp"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/collision"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/movement"
	"github.com/krishnaZawar/LevelCraft/utils/input"
	"github.com/stretchr/testify/assert"
)

var inputMapping = &input.InputMapping{}

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
		{
			comp:         movement.NewBaseInputMovement(inputMapping),
			expectedName: base.ComponentName_InputMovement,
		},
		{
			comp:         collision.NewBaseCollider(),
			expectedName: base.ComponentName_Collider,
		},
	}

	for _, tt := range tests {
		name := tt.comp.GetComponentName()
		assert.Equal(t, tt.expectedName, name)
	}
}

func Test_ComponentList(t *testing.T) {
	listLen := 4
	assert.Equal(t, listLen, len(ComponentList))

	assert.Equal(t, true, contains(ComponentList, base.ComponentName_InputMovement))
	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Transform))
	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Color))
	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Collider))
}
