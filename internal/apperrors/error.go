// Package apperrors implements the application errors.
package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrURLNotFound                  = errors.New("url not found")
	ErrUnableToGetUserIDFromContext = errors.New("unable to get user id from context")
	ErrEmptyRequest                 = errors.New("unable to handle empty request")
	ErrDuplicatedKeys               = errors.New("duplicated keys in batch")
	ErrURLAlreadyExists             = errors.New("url already exists")
	ErrURLDeleted                   = errors.New("url deleted")
)

type ValueError struct {
	caller  string
	message string
	err     error
}

// NewValueError creates a new ValueError with the given message, caller, and error.
//
// Parameters:
//
//	message string - the error message
//	caller string - the caller of the function tracing place in the code
//	err error - the original error
//
// Return type:
//
//	error - the newly created ValueError
func NewValueError(message string, caller string, err error) error {
	return &ValueError{
		caller:  caller,
		message: message,
		err:     err,
	}
}

// Error returns a string representing the error.
//
// No parameters.
// Returns a string.
func (v *ValueError) Error() string {
	return fmt.Sprintf("%s %s %s", v.caller, v.message, v.err)
}

// Unwrap returns the error that has been wrapped by ValueError.
// No parameters. Returns an error.
func (v *ValueError) Unwrap() error {
	return v.err
}
