package base

import (
	"errors"
	"fmt"
)

// Rejected rather than clamped: a negative speed is a configuration mistake
// that would silently invert the direction of travel.
var ErrNegativeSpeed = errors.New("speed cannot be negative")

// A negative stop distance would mean the object has to travel past its
// target before it counts as arrived.
var ErrNegativeStopDistance = errors.New("stop distance cannot be negative")

func NewErrUnknownFollowTarget(targetMode string) error {
	return fmt.Errorf("unknown follow target mode: %q", targetMode)
}
