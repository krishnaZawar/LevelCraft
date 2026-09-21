package component

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
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

type MockComponent struct {
	mockGetComponentName    func() string
	mockGetComponentDetails func() ComponentDetails
	mockBuildFromDetails    func(json.RawMessage) error
}

func (mc *MockComponent) GetComponentName() string {
	return mc.mockGetComponentName()
}
func (mc *MockComponent) GetComponentDetails() ComponentDetails {
	return mc.mockGetComponentDetails()
}
func (mc *MockComponent) BuildFromDetails(data json.RawMessage) error {
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
			mockGetComponentDetails: func() ComponentDetails {
				return ComponentDetails{
					Name: componentName,
					Data: json.RawMessage([]byte{}),
				}
			},
			mockBuildFromDetails: func(m json.RawMessage) error {
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
