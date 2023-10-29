package apperrors

import "errors"

var (
	ErrorUrlNotFound error = errors.New("url not found")
	ErrorEmptyRequest error = errors.New("unable to handle empty request")
)

