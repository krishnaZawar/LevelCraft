package component

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

type MockComponent struct {
	mockGetComponentName    func() string
	mockGetComponentDetails func() map[string]interface{}
	mockBuildFromDetails    func(map[string]interface{})
}

func (mc *MockComponent) GetComponentName() string {
	return mc.mockGetComponentName()
}
func (mc *MockComponent) GetComponentDetails() map[string]interface{} {
	return mc.mockGetComponentDetails()
}
func (mc *MockComponent) BuildFromDetails(data map[string]interface{}) {
	mc.mockBuildFromDetails(data)
}

func Test_RegisterAndFetch(t *testing.T) {
	compRegistry := newComponentRegistry()

	const (
		componentName = "test-component"
	)

	var (
		comp = &MockComponent{
			mockGetComponentName: func() string {
				return componentName
			},
			mockGetComponentDetails: func() map[string]interface{} {
				return map[string]interface{}{}
			},
			mockBuildFromDetails: func(m map[string]interface{}) {
				// pass
			},
		}
	)

	compRegistry.register(componentName, comp)

	component, ok := compRegistry.GetComponent(componentName)

	assert.Equal(t, true, ok)
	assert.Equal(t, comp, component)
}

func Test_NewComponentRegistry(t *testing.T) {
	compRegistry := NewComponentRegistry()

	var expected Component
	comp, ok := compRegistry.GetComponent(base.ComponentName_Transform)
	expected = newBaseTransform()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)

	comp, ok = compRegistry.GetComponent(base.ComponentName_Color)
	expected = newBaseColor()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)
}

func Test_GetComponentName(t *testing.T) {
	tests := []struct {
		comp         Component
		expectedName string
	}{
		{
			comp:         newBaseTransform(),
			expectedName: base.ComponentName_Transform,
		},
		{
			comp:         newBaseColor(),
			expectedName: base.ComponentName_Color,
		},
	}

	for _, tt := range tests {
		name := tt.comp.GetComponentName()
		assert.Equal(t, tt.expectedName, name)
	}
}

func Test_GetComponentDetails(t *testing.T) {
	tests := []struct {
		comp            Component
		expectedDetails map[string]interface{}
	}{
		{
			comp: NewTransform(100, 100, 100, 100),
			expectedDetails: map[string]interface{}{
				"x": 100, "y": 100, "w": 100, "h": 100,
			},
		},
		{
			comp: NewColor(0, 0, 0, 255),
			expectedDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 0, "a": 255,
			},
		},
	}

	for _, tt := range tests {
		details := tt.comp.GetComponentDetails()
		assert.Equal(t, tt.expectedDetails, details)
	}
}

func Test_BuildFromDetails(t *testing.T) {
	tests := []struct {
		comp         Component
		buildDetails map[string]interface{}
	}{
		{
			comp: newBaseTransform(),
			buildDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
		},
		{
			comp: newBaseColor(),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 0, "a": 0,
			},
		},
	}

	for _, tt := range tests {
		tt.comp.BuildFromDetails(tt.buildDetails)
		details := tt.comp.GetComponentDetails()
		assert.Equal(t, tt.buildDetails, details)
	}
}
