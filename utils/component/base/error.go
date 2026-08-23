package base

import (
	"errors"
	"fmt"
)

var (
	ErrExpectedInteger = errors.New("ComponentError: Expected integer value")
)

var (
	ErrColorValueRangeOutOfBounds = fmt.Errorf("ComponentError: Color values should be in the range of %d to %d", ColorValueRangeMin, ColorValueRangeMax)
)
