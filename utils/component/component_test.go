package component

import (
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

type MockComponent struct {
	mockGetComponentName    func() string
	mockGetComponentDetails func() map[string]interface{}
	mockBuildFromDetails    func(map[string]interface{}) error
}

func (mc *MockComponent) GetComponentName() string {
	return mc.mockGetComponentName()
}
func (mc *MockComponent) GetComponentDetails() map[string]interface{} {
	return mc.mockGetComponentDetails()
}
func (mc *MockComponent) BuildFromDetails(data map[string]interface{}) error {
	return mc.mockBuildFromDetails(data)
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
			mockBuildFromDetails: func(m map[string]interface{}) error {
				return nil
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
		name                string
		comp                Component
		buildDetails        map[string]interface{}
		expectedDetails     map[string]interface{}
		expectedReturnValue error
	}{
		{
			name: "TransformTest - all values correct",
			comp: newBaseTransform(),
			buildDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
			expectedDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
			expectedReturnValue: nil,
		},
		{
			name: "TransformTest - x is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"x": "a", "y": 0, "w": 0, "h": 0,
			},
			expectedDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "TransformTest - y is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"x": 0, "y": "a", "w": 0, "h": 0,
			},
			expectedDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "TransformTest - w is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": "w", "h": 0,
			},
			expectedDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "TransformTest - h is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": "h",
			},
			expectedDetails: map[string]interface{}{
				"x": 0, "y": 0, "w": 0, "h": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - all values correct",
			comp: newBaseColor(),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: nil,
		},
		{
			name: "ColorTest - r is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": "0", "g": 0, "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - r is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": -1, "g": 0, "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - r is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 256, "g": 0, "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - g is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": "0", "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - g is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": -1, "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - g is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 256, "b": 0, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - b is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": "0", "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - b is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": -1, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - b is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 256, "a": 0,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - a is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 0, "a": "0",
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - a is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 0, "a": -1,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - a is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				"r": 0, "g": 0, "b": 255, "a": 256,
			},
			expectedDetails: map[string]interface{}{
				"r": 10, "g": 0, "b": 0, "a": 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
	}

	for _, tt := range tests {
		testName := map[string]interface{}{
			"name": tt.name,
		}
		err := tt.comp.BuildFromDetails(tt.buildDetails)
		assert.Equal(t, tt.expectedReturnValue, err, testName)
		details := tt.comp.GetComponentDetails()
		assert.Equal(t, tt.expectedDetails, details, testName)
	}
}
