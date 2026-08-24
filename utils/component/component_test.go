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

func Test_ComponentList(t *testing.T) {
	listLen := 2
	assert.Equal(t, listLen, len(ComponentList))

	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Transform))
	assert.Equal(t, true, contains(ComponentList, base.ComponentName_Color))
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
				Transform_CurLabelX: 100, Transform_CurLabelY: 100, Transform_CurLabelW: 100, Transform_CurLabelH: 100,
			},
		},
		{
			comp: NewColor(0, 0, 0, 255),
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 255,
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
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedReturnValue: nil,
		},
		{
			name: "TransformTest - x is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Transform_CurLabelX: "a", Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "TransformTest - y is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: "a", Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "TransformTest - w is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: "a", Transform_CurLabelH: 0,
			},
			expectedDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "TransformTest - h is not integer",
			comp: NewTransform(0, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: "a",
			},
			expectedDetails: map[string]interface{}{
				Transform_CurLabelX: 0, Transform_CurLabelY: 0, Transform_CurLabelW: 0, Transform_CurLabelH: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - all values correct",
			comp: newBaseColor(),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: nil,
		},
		{
			name: "ColorTest - r is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: "0", Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - r is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: -1, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - r is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 256, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - g is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: "0", Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - g is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: -1, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - g is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 256, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - b is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: "0", Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - b is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: -1, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - b is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 256, Color_CurLabelA: 0,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - a is not integer",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: "0",
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrExpectedInteger,
		},
		{
			name: "ColorTest - a is below range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: -1,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
			},
			expectedReturnValue: base.ErrColorValueRangeOutOfBounds,
		},
		{
			name: "ColorTest - a is above range",
			comp: NewColor(10, 0, 0, 0),
			buildDetails: map[string]interface{}{
				Color_CurLabelR: 0, Color_CurLabelG: 0, Color_CurLabelB: 255, Color_CurLabelA: 256,
			},
			expectedDetails: map[string]interface{}{
				Color_CurLabelR: 10, Color_CurLabelG: 0, Color_CurLabelB: 0, Color_CurLabelA: 0,
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
