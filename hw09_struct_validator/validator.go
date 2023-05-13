package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct, got: %v", value.Kind())
	}

	var validationErrors ValidationErrors
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldValue := value.Field(i)
		fieldErrors := validateField(field.Name, fieldValue, tag)
		validationErrors = append(validationErrors, fieldErrors...)
	}

	if len(validationErrors) == 0 {
		return nil
	}
	return validationErrors
}

func validateField(fieldName string, fieldValue reflect.Value, tag string) []ValidationError {
	var validationErrors []ValidationError

	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		ruleParts := strings.SplitN(rule, ":", 2)
		if len(ruleParts) != 2 {
			validationErrors = append(validationErrors, ValidationError{fieldName, fmt.Errorf("invalid rule format: %s", rule)})
			continue
		}

		validatorName := ruleParts[0]
		validatorParam := ruleParts[1]

		switch fieldValue.Kind() { //nolint: exhaustive // we don't need to check all types - we have default case
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := validateInt(fieldName, fieldValue.Int(), validatorName, validatorParam)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{fieldName, err})
			}
		case reflect.String:
			err := validateString(fieldName, fieldValue.String(), validatorName, validatorParam)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{fieldName, err})
			}
		case reflect.Slice:
			for i := 0; i < fieldValue.Len(); i++ {
				elemErrors := validateField(fmt.Sprintf("%s[%d]", fieldName, i), fieldValue.Index(i), tag)
				validationErrors = append(validationErrors, elemErrors...)
			}
		default:
			err := UnsupportedFieldTypeError{Field: fieldName, Type: fieldValue.Kind().String()}
			validationErrors = append(validationErrors, ValidationError{fieldName, err})
		}
	}

	return validationErrors
}

func validateInt(field string, value int64, validatorName, validatorParam string) error {
	switch validatorName {
	case "min":
		min, err := strconv.ParseInt(validatorParam, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid min validator parameter: %w", err)
		}
		if value < min {
			return ErrValidationMin{Field: field, Min: int(min)}
		}
	case "max":
		max, err := strconv.ParseInt(validatorParam, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max validator parameter: %w", err)
		}
		if value > max {
			return ErrValidationMax{Field: field, Max: int(max)}
		}
	case "in":
		inValues := strings.Split(validatorParam, ",")
		found := false
		for _, inValue := range inValues {
			inInt, err := strconv.ParseInt(inValue, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid in validator parameter: %w", err)
			}
			if value == inInt {
				found = true
				break
			}
		}
		if !found {
			return ErrValidationIn{Field: field, Set: inValues}
		}
	default:
		return ErrUnsupportedValidator{Field: field, Validator: validatorName, TargetType: "int"}
	}
	return nil
}

func validateString(field, value, validatorName, validatorParam string) error {
	switch validatorName {
	case "len":
		length, err := strconv.Atoi(validatorParam)
		if err != nil {
			return fmt.Errorf("invalid len validator parameter: %w", err)
		}
		if len(value) != length {
			return ErrValidationLen{Field: field, Len: length}
		}
	case "regexp":
		pattern, err := regexp.Compile(validatorParam)
		if err != nil {
			return fmt.Errorf("invalid regexp validator parameter: %w", err)
		}
		if !pattern.MatchString(value) {
			return ErrValidationRegexp{Field: field, Pattern: pattern.String()}
		}
	case "in":
		inValues := strings.Split(validatorParam, ",")
		found := false
		for _, inValue := range inValues {
			if value == inValue {
				found = true
				break
			}
		}
		if !found {
			return ErrValidationIn{Field: field, Set: inValues}
		}
	default:
		return ErrUnsupportedValidator{Field: field, Validator: validatorName, TargetType: "string"}
	}

	return nil
}
