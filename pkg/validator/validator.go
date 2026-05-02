package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func TranslateValidationErrors(err error) []FieldError {
	var ve validator.ValidationErrors
	if ok := isValidationErrors(err, &ve); !ok {
		return nil
	}

	errs := make([]FieldError, 0, len(ve))
	for _, fe := range ve {
		errs = append(errs, FieldError{
			Field:   toSnakeCase(fe.Field()),
			Message: translateTag(fe),
		})
	}
	return errs
}

func isValidationErrors(err error, target *validator.ValidationErrors) bool {
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		*target = ve
	}
	return ok
}

func translateTag(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		if fe.Type().Kind().String() == "string" {
			return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		if fe.Type().Kind().String() == "string" {
			return fmt.Sprintf("%s must be at most %s characters long", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at most %s", field, fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be a valid number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, fe.Param())
	case "nefield":
		return fmt.Sprintf("%s must not match %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s failed validation on '%s' rule", field, fe.Tag())
	}
}

func toSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	result := make([]rune, 0, len(s)+4)
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, r+32)
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
