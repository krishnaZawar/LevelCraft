package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockComponent struct {
	mockGetComponentName    func() string
	mockGetComponentDetails func() ComponentDetails
	mockBuildFromDetails    func(json.RawMessage) error
	mockHandlesInput        func(string) bool
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
func (mc *MockComponent) HandlesInput(input string) bool {
	return mc.mockHandlesInput(input)
}

func Test_RegisterAndFetch(t *testing.T) {
	compRegistry := NewComponentRegistry()

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
			mockHandlesInput: func(input string) bool {
				return false
			},
		}

		compFactory = func() Component {
			return comp
		}
	)

	compRegistry.Register(componentName, compFactory)

	factory, ok := compRegistry.GetComponent(componentName)
	assert.Equal(t, true, ok)

	component := factory()
	assert.Equal(t, comp, component)
}
