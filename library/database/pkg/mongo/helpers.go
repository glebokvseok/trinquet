package mongo

import "github.com/pkg/errors"

const (
	errorMessage = "unhandled panic occurred while executing mongodb request: %v"
)

func PanicHandler(returnErr *error) {
	if r := recover(); r != nil {
		if returnErr != nil {
			*returnErr = errors.Errorf(errorMessage, r)
		}
	}
}
