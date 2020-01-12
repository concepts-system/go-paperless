package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type genericError struct {
	cause   error
	root    error
	context map[string]string
}

func (err *genericError) Error() string {
	return err.cause.Error()
}

// New constructs a new error with the given message.
func New(message string) error {
	cause := errors.New(message)

	return &genericError{
		cause: cause,
		root:  cause,
	}
}

// Newf constructs a new error with the given message format and args.
func Newf(message string, args ...interface{}) error {
	cause := errors.New(fmt.Sprintf(message, args...))

	return &genericError{
		cause: cause,
		root:  cause,
	}
}

// Wrap wraps the given error with a new message.
func Wrap(err error, message string) error {
	return Wrapf(err, message)
}

// Wrapf wraps the given error with a new message format and arguments.
func Wrapf(err error, message string, args ...interface{}) error {
	wrappedErr := &genericError{
		cause: errors.Wrapf(err, message, args...),
	}

	if genericErr, ok := err.(*genericError); ok {
		wrappedErr.root = genericErr.root
		wrappedErr.context = genericErr.context
	} else {
		wrappedErr.root = wrappedErr.cause
	}

	return wrappedErr
}

// RootCause returns the root cause for the given error.
// Might be the error itself in case not root cause could be determined.
func RootCause(err error) error {
	if genericErr, ok := err.(*genericError); ok {
		return genericErr.root
	}

	return err
}

// AddContext adds context to the given error.
func AddContext(err error, field, message string) error {
	if genericErr, ok := err.(*genericError); ok {
		if genericErr.context == nil {
			genericErr.context = map[string]string{}
		}

		genericErr.context[field] = message
		return genericErr
	}

	return &genericError{
		cause: err,
		root:  err,
		context: map[string]string{
			field: message,
		},
	}
}

// GetContext returns an error's context.
func GetContext(err error) map[string]string {
	if customErr, ok := err.(*genericError); ok && len(customErr.context) > 0 {
		return customErr.context
	}

	return nil
}
