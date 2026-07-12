package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRegistry(t *testing.T) {
	registry := NewRegistry[string, string]()
	assert.NotNil(t, registry.items)
}

func Test_Register(t *testing.T) {
	tests := []struct {
		registry *Registry[string, string]
		key      string
		value    string
	}{
		{
			registry: NewRegistry[string, string](),
			key:      "test-key",
			value:    "test-value",
		},
		{
			registry: &Registry[string, string]{},
			key:      "test-key",
			value:    "test-value",
		},
	}

	for _, tt := range tests {
		tt.registry.Register(tt.key, tt.value)
		value, ok := tt.registry.GetValue(tt.key)
		if !ok {
			t.Errorf("Value not found for key %s", tt.key)
		}
		assert.Equal(t, true, ok)
		assert.Equal(t, tt.value, value)
	}
}

func Test_GetValue(t *testing.T) {
	registry := &Registry[string, string]{}
	key := "test-key"
	value := "test-value"
	t.Run("Test Get before registration", func(t *testing.T) {
		_, ok := registry.GetValue(key)
		assert.Equal(t, false, ok)
	})

	t.Run("Test Get after registration", func(t *testing.T) {
		registry.Register(key, value)
		val, ok := registry.GetValue(key)
		assert.Equal(t, true, ok)
		assert.Equal(t, value, val)
	})
}
