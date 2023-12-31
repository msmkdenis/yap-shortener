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

func NewValueError(message string, caller string, err error) error {
	return &ValueError{
		caller:  caller,
		message: message,
		err:     err,
	}
}

func (v *ValueError) Error() string {
	return fmt.Sprintf("%s %s %s", v.caller, v.message, v.err)
}

func (v *ValueError) Unwrap() error {
	return v.err
}
