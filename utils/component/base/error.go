package base

import (
	"errors"
	"fmt"
)

var (
	ErrExpectedInteger = errors.New("error: Expected integer value")
)

var (
	ErrColorValueRangeOutOfBounds = fmt.Errorf("error: Color values should be in the range of %d to %d", ColorValueRangeMin, ColorValueRangeMax)
)
