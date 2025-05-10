package core

import (
	"math"
	"strconv"
	"strings"
)

type InputValidator struct{}

func (v InputValidator) ParseInt(input string) (int, error) {
	input = strings.TrimSpace(input)
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0, NewAppError(ErrInvalidInput, "invalid_number_input", nil)
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, NewAppError(ErrOutOfRange, "out_of_range_generic", nil)
	}
	return num, nil
}

func (v InputValidator) IsNumInRange(num, min, max int) (bool, error) {
	if num > math.MaxInt32 || num < math.MinInt32 {
		return false, NewAppError(ErrOutOfRange, "out_of_range_generic", nil)
	}
	if num < min {
		return false, NewAppError(ErrOutOfRange, "out_of_range_min", map[string]any{
			"min": min,
		})
	}
	if num > max {
		return false, NewAppError(ErrOutOfRange, "out_of_range_max", map[string]any{
			"max": max,
		})
	}
	return true, nil
}

func (v InputValidator) ParseIntInRange(input string, min, max int) (int, error) {
	num, err := v.ParseInt(input)
	if err != nil {
		return 0, err
	}
	if _, err := v.IsNumInRange(num, min, max); err != nil {
		return 0, err
	}
	return num, nil
}

func (v InputValidator) ParseOptionalIntInRange(input string, defaultValue, min, max int) (int, error) {
	if input == "" {
		if _, err := v.IsNumInRange(defaultValue, min, max); err != nil {
			return 0, err
		}
		return defaultValue, nil
	}
	return v.ParseIntInRange(input, min, max)
}
