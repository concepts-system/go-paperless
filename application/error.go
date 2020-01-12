package application

import (
	"github.com/concepts-system/go-paperless/errors"
)

const fieldErrorType = "__errorType"

// ErrorType enumerates all possible application error types.
type ErrorType string

const (
	// InternalServerError specifies an unknown error.
	InternalServerError = ErrorType(iota)
	// BadRequestError specifies a validation realted error.
	BadRequestError
	// UnauthorizedError specifies authentication related errors.
	UnauthorizedError
	// ForbiddenError specifies authorization (permission) related errors.
	ForbiddenError
	// NotFoundError specifies errors related with non-existent resources.
	NotFoundError
	// ConflictError specifies errors related with a resource conflict.
	ConflictError
	// UnexpectedError specifies errors occurring unexpectedly, caused by technical issues.
	UnexpectedError
)

// New creates a new error for the given message.
func (typ ErrorType) New(message string) error {
	return typ.Newf(message)
}

// Newf creates a new error with the given message format and arguments.
func (typ ErrorType) Newf(message string, args ...interface{}) error {
	return SetErrorType(errors.Newf(message, args...), typ)
}

// SetErrorType associates the given error type with the given error.
func SetErrorType(err error, typ ErrorType) error {
	return errors.AddContext(err, fieldErrorType, string(typ))
}

// RemoveErrorType removes associated error type information form the given error.
func RemoveErrorType(err error) error {
	context := errors.GetContext(err)

	if context != nil {
		delete(context, fieldErrorType)
	}

	return err
}

// GetErrorType retreives an error's associated error type.
func GetErrorType(err error) ErrorType {
	context := errors.GetContext(err)
	if context == nil {
		return InternalServerError
	}

	if typ, ok := context[fieldErrorType]; ok {
		return ErrorType(typ)
	}

	return InternalServerError
}
