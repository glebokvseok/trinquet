package errors

import (
	"fmt"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"net/http"
)

const (
	UserNotFoundError      = "USER_NOT_FOUND_ERROR"
	IncorrectPasswordError = "INCORRECT_PASSWORD_ERROR"
	UserAlreadyExistsError = "USER_ALREADY_EXISTS_ERROR"
)

func NewUserNotFoundError(email string) *errors.APIError {
	return errors.NewAPIError(
		http.StatusNotFound,
		UserNotFoundError,
		fmt.Sprintf("user with email: %s has not been registered yet.", email),
	)
}

func NewIncorrectPasswordError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusUnauthorized,
		IncorrectPasswordError,
		"Incorrect password.",
	)
}

func NewUserAlreadyExistsError(email string) *errors.APIError {
	return errors.NewAPIError(
		http.StatusConflict,
		UserAlreadyExistsError,
		fmt.Sprintf("user with email: %s has been registered already.", email),
	)
}
