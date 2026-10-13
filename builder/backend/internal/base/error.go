package base

import (
	"errors"
	"fmt"
)

var (
	NewErrComponentTypeMisMatch = func(expected, found string) error {
		return fmt.Errorf("componentMisMatchError: expected component %s but found %s", expected, found)
	}

	NewErrGameobjectNotFound = func(gameobjectID string) error {
		return fmt.Errorf("error: Gameobject with id %s not found in the game scene", gameobjectID)
	}

	NewErrEventTypeMismatch = func(expected string, found string) error {
		return fmt.Errorf("eventMismatchError: expected event %s but found event %s", expected, found)
	}
	NewErrComponentNotFound = func(compName string, gameobjectID string) error {
		return fmt.Errorf("error: gameobject with id %s does not have component %s attached to it", gameobjectID, compName)
	}

	ErrNilComponentPassed = errors.New("eventGeneratorErr: component was not passed in the context")
)

var (
	MsgGameobjectNotFound        = "gameobject not found"
	MsgExpectedEventTypeNotFound = "expected event type not found"
	MsgComponentNotFound         = "component not found"
)
