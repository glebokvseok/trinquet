package errors

import (
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"net/http"
)

const (
	MalformedOrExpiredSignatureError = "INVALID_OR_EXPIRED_JWT_ERROR"
)

func NewInvalidOrExpiredJWTError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusUnauthorized,
		MalformedOrExpiredSignatureError,
		"Invalid or expired jwt.",
	)
}
