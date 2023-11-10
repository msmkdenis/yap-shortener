package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrorURLNotFound  error = errors.New("url not found")
	ErrorEmptyRequest error = errors.New("unable to handle empty request")
	ErrorDuplicatedKeys error = errors.New("duplicated keys in batch")
	ErrorURLAlreadyExists error = errors.New("url already exists")
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
	return fmt.Sprintf("caller: %s message: %s error: %s", v.caller, v.message, v.err)
}

func (v *ValueError) Unwrap() error {
	return v.err
}
