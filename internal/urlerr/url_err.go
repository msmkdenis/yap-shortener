// Package url_err contains the errors for the url domain.
package urlerr

import "errors"

// Errors
var (
	ErrURLNotFound                  = errors.New("url not found")
	ErrURLDeleted                   = errors.New("url deleted")
	ErrUnableToGetUserIDFromContext = errors.New("unable to get user id from context")
	ErrEmptyRequest                 = errors.New("unable to handle empty request")
	ErrDuplicatedKeys               = errors.New("duplicated keys in batch")
	ErrURLAlreadyExists             = errors.New("url already exists")
)
