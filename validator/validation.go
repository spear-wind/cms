package validator

import (
	"regexp"
	"strings"
)

type ValidationError struct {
	FieldName    string `json:"field_name"`
	ErrorMessage string `json:"error_message"`
}

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) AddError(fieldName string, errorMessage string) {
	r.Errors = append(r.Errors, NewValidationError(fieldName, errorMessage))
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) != 0
}

func NewValidationError(fieldName string, errorMessage string) ValidationError {
	return ValidationError{
		FieldName:    fieldName,
		ErrorMessage: errorMessage,
	}
}

func NewValidationResult() ValidationResult {
	return ValidationResult{
		Errors: []ValidationError{},
	}
}

func ValidateEmail(email string) bool {
	email = strings.ToLower(email)
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}
