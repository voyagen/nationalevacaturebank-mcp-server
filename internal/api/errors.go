package api

import "fmt"

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int
	Message    string
	Endpoint   string
	Cause      error
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("API error %d at %s: %s (caused by: %v)",
			e.StatusCode, e.Endpoint, e.Message, e.Cause)
	}
	return fmt.Sprintf("API error %d at %s: %s", e.StatusCode, e.Endpoint, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Cause
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, message, endpoint string, cause error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Endpoint:   endpoint,
		Cause:      cause,
	}
}

// ValidationError represents a parameter validation error
type ValidationError struct {
	Parameter string
	Value     interface{}
	Message   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for parameter '%s' (value: %v): %s",
		e.Parameter, e.Value, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(parameter string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Parameter: parameter,
		Value:     value,
		Message:   message,
	}
}