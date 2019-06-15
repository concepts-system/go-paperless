package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorType defines an enum for all possible kinds of errors.
type ErrorType string

/* Custom Error */

type customError struct {
	typ     ErrorType
	cause   error
	context map[string]string
}

const (
	// Unknown specifies an unknown error.
	Unknown = ErrorType(iota)
	// BadRequest specifies a validation realted error.
	BadRequest
	// Unauthorized specifies authentication related errors.
	Unauthorized
	// Forbidden specifies authorization (permission) related errors.
	Forbidden
	// NotFound specifies errors related with non-existent resources.
	NotFound
	// Conflict specifies errors related with a resource conflict.
	Conflict
	// Unexpected specifies errors occurring unexpectedly, caused by technical issues.
	Unexpected
)

func (err customError) Error() string {
	return err.cause.Error()
}

// New creates a new error for the given message.
func (typ ErrorType) New(message string) error {
	return customError{
		typ:   typ,
		cause: errors.New(message),
	}
}

// Newf creates a new error with the given message format and arguments.
func (typ ErrorType) Newf(message string, args ...interface{}) error {
	return customError{
		typ:   typ,
		cause: fmt.Errorf(message, args...),
	}
}

// Wrap wraps a given error with a new message.
func (typ ErrorType) Wrap(err error, message string) error {
	return typ.Wrapf(err, message)
}

// Wrapf wraps a given error with a new message format and arguments.
func (typ ErrorType) Wrapf(err error, message string, args ...interface{}) error {
	return customError{
		typ:   typ,
		cause: errors.Wrapf(err, message, args...),
	}
}

/* Public Error Construction Methods */

// New constructs a new error with the given message.
func New(message string) error {
	return customError{
		typ:   Unknown,
		cause: errors.New(message),
	}
}

// Newf constructs a new error with the given message format and args.
func Newf(message string, args ...interface{}) error {
	return customError{
		typ:   Unknown,
		cause: errors.New(fmt.Sprintf(message, args...)),
	}
}

// Wrap wraps the given error with a new message.
func Wrap(err error, message string) error {
	return Wrapf(err, message)
}

// Wrapf wraps the given error with a new message format and arguments.
func Wrapf(err error, message string, args ...interface{}) error {
	wrappedErr := errors.Wrapf(err, message, args...)

	if customErr, ok := err.(customError); ok {
		return customError{
			typ:     customErr.typ,
			cause:   wrappedErr,
			context: customErr.context,
		}
	}

	return customError{
		typ:   Unknown,
		cause: wrappedErr,
	}
}

// Cause returns the cause of the given error.
func Cause(err error) error {
	return errors.Cause(err)
}

// GetType returns the type of the given error.
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.typ
	}

	return Unknown
}

/* Custom Error Context */

// AddContext adds context to the given error.
func AddContext(err error, field, message string) error {
	if customErr, ok := err.(customError); ok {
		if customErr.context == nil {
			customErr.context = map[string]string{}
		}

		customErr.context[field] = message
		return customErr
	}

	return customError{
		typ:   Unknown,
		cause: err,
		context: map[string]string{
			field: message,
		},
	}
}

// GetContext returns an error's context.
func GetContext(err error) *map[string]string {
	if customErr, ok := err.(customError); ok && len(customErr.context) > 0 {
		return &customErr.context
	}

	return nil
}
