package hw09structvalidator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStrings := make([]string, 0, len(v))

	for _, err := range v {
		errStrings = append(errStrings, fmt.Sprintf("%s: %s", err.Field, err.Err.Error()))
	}
	return strings.Join(errStrings, "; ")
}

type ErrValidationRegexp struct {
	Field   string
	Pattern string
}

func (e ErrValidationRegexp) Error() string {
	return fmt.Sprintf("%s: value must match pattern %s", e.Field, e.Pattern)
}

type ErrUnsupportedValidator struct {
	Field      string
	Validator  string
	TargetType string
}

func (e ErrUnsupportedValidator) Error() string {
	return fmt.Sprintf("unsupported validator for %s: %s in field %s", e.TargetType, e.Validator, e.Field)
}

type UnsupportedFieldTypeError struct {
	Field string
	Type  string
}

func (e UnsupportedFieldTypeError) Error() string {
	return fmt.Sprintf("unsupported field type: %s for field: %s", e.Type, e.Field)
}

type ErrValidationLen struct {
	Field string
	Len   int
}

func (e ErrValidationLen) Error() string {
	return fmt.Sprintf("%s: value must have a length of %d", e.Field, e.Len)
}

type ErrValidationMin struct {
	Field string
	Min   int
}

func (e ErrValidationMin) Error() string {
	return fmt.Sprintf("%s: value must be at least %d", e.Field, e.Min)
}

type ErrValidationMax struct {
	Field string
	Max   int
}

func (e ErrValidationMax) Error() string {
	return fmt.Sprintf("%s: value must be no more than %d", e.Field, e.Max)
}

type ErrValidationIn struct {
	Field string
	Set   []string
}

func (e ErrValidationIn) Error() string {
	return fmt.Sprintf("%s: value must be in set %v", e.Field, e.Set)
}
