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

	comp, ok := compRegistry.GetComponent(base.ComponentName_Transform)
	expected := newBaseTransform()
	assert.Equal(t, true, ok)
	assert.Equal(t, expected, comp)
}
