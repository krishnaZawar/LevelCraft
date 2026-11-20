package controlmanager

import "fmt"

// represents the current execution state of the game loop
type ExecutionState string

const (
	ExecutionState_Inactive ExecutionState = "inactive"
	ExecutionState_Active   ExecutionState = "active"
)

// checks whether the executionState passed is a valid one or not
func (e ExecutionState) isValid() bool {
	switch e {
	case ExecutionState_Active, ExecutionState_Inactive:
		return true
	}
	return false
}

// ControlManager handles access to the execution state and tells whether the game loop should execute or not
type ControlManager struct {
	// execution state of the game loop
	// default value should be inactive
	ExecutionState ExecutionState
}

func newControlManager() *ControlManager {
	return &ControlManager{
		ExecutionState: ExecutionState_Inactive,
	}
}

// validates and sets the execution state of the control manager
func (cm *ControlManager) SetExecutionState(ExecutionState ExecutionState) error {
	if !ExecutionState.isValid() {
		return fmt.Errorf("error: Invalid ExecutionState value passed: %s", ExecutionState)
	}
	cm.ExecutionState = ExecutionState
	return nil
}

// returns the current execution state held by the control manager
func (cm *ControlManager) GetExecutionState() ExecutionState {
	return cm.ExecutionState
}

// tells whether the game loop should process a frame or not
func (cm *ControlManager) CanExecuteFrame() bool {
	return cm.ExecutionState == ExecutionState_Active
}

var controlManager = newControlManager()

func Get() *ControlManager {
	return controlManager
}
