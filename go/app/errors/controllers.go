// app/errors/controllers.go

package errors

import "errors"

var (
	ErrLimitViolation = errors.New("the limit value violates the server threshold")
)
