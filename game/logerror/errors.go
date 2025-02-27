package logerror

import (
	"errors"
	"fmt"
	"strings"
)

// Error represents a custom error type that ensures capitalized error messages
type LogError struct {
	err error
}

// New creates a new Error with a capitalized message
func New(text string) error {
	if text == "" {
		return nil
	}

	return &LogError{
		err: errors.New(text),
	}
}

// Errorf creates a new Error with a formatted message
func Errorf(format string, args ...interface{}) error {
	return &LogError{
		err: fmt.Errorf(capitalize(format), args...),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &LogError{
		err: fmt.Errorf("%s: %w", capitalize(message), err),
	}
}

// Error implements the error interface
func (e *LogError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}
	return e.err.Error()
}

// Unwrap returns the wrapped error
func (e *LogError) Unwrap() error {
	return errors.Unwrap(e.err)
}

// Is reports whether any error in err's chain matches target
func (e *LogError) Is(target error) bool {
	return errors.Is(e.err, target)
}

// As finds the first error in err's chain that matches target
func (e *LogError) As(target interface{}) bool {
	return errors.As(e.err, target)
}

// capitalize ensures the first letter of the string is uppercase
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Common errors that can be reused throughout the application
var (
	ErrInvalidInput   = New("invalid input")
	ErrNotFound       = New("not found")
	ErrUnauthorized   = New("unauthorized")
	ErrInternal       = New("internal error")
	ErrNotImplemented = New("not implemented")
)
