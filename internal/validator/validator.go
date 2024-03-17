package validator

import (
	"slices"
	"strconv"
	"strings"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

// validate Name is not empty
// validate Rating is float64 from 0.0 to 5.0
// validate NumReviews is positive int

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, ok := v.FieldErrors[key]; !ok {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key string, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func PositiveInt(value string) bool {
	num, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return num > 0
}

func IsFloat64(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func InRangeFloat64(value, min, max float64) bool {
	return value >= min && value <= max
}
