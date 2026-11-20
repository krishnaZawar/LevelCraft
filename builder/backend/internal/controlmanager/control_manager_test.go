package controlmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidExecutionState(t *testing.T) {
	tests := []struct {
		name  string
		value ExecutionState
		valid bool
	}{
		{
			name:  "ExecutionState active",
			value: ExecutionState_Active,
			valid: true,
		},
		{
			name:  "ExecutionState inactive",
			value: ExecutionState_Inactive,
			valid: true,
		},
		{
			name:  "invalid ExecutionState",
			value: ExecutionState("invalid ExecutionState"),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.value.isValid())
		})
	}
}

func Test_NewControlManager(t *testing.T) {
	manager := newControlManager()

	assert.Equal(t, ExecutionState_Inactive, manager.ExecutionState)
}

func Test_SetExecutionState(t *testing.T) {
	t.Run("invalid ExecutionState", func(t *testing.T) {
		manager := &ControlManager{
			ExecutionState: ExecutionState_Inactive,
		}

		err := manager.SetExecutionState(ExecutionState("invalid ExecutionState"))
		assert.NotNil(t, err)
		assert.Equal(t, ExecutionState_Inactive, manager.ExecutionState)
	})

	t.Run("valid ExecutionState", func(t *testing.T) {
		manager := &ControlManager{
			ExecutionState: ExecutionState_Inactive,
		}

		err := manager.SetExecutionState(ExecutionState_Active)
		assert.Nil(t, err)
		assert.Equal(t, ExecutionState_Active, manager.ExecutionState)
	})
}

func Test_GetExecutionState(t *testing.T) {
	expectedType := ExecutionState("type")
	manager := &ControlManager{
		ExecutionState: expectedType,
	}

	assert.Equal(t, expectedType, manager.GetExecutionState())
}

func Test_CanExecuteFrame(t *testing.T) {
	t.Run("cannot execute frame", func(t *testing.T) {
		manager := &ControlManager{
			ExecutionState: ExecutionState_Inactive,
		}
		assert.Equal(t, false, manager.CanExecuteFrame())
	})

	t.Run("can execute frame", func(t *testing.T) {
		manager := &ControlManager{
			ExecutionState: ExecutionState_Active,
		}
		assert.Equal(t, true, manager.CanExecuteFrame())
	})
}
