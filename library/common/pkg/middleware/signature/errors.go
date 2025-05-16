package signmw

import (
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"net/http"
)

const (
	MalformedOrExpiredSignatureError = "MALFORMED_OR_EXPIRED_SIGNATURE_ERROR"
)

func NewMalformedOrExpiredSignatureError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusUnauthorized,
		MalformedOrExpiredSignatureError,
		"Request signature is malformed or expired.",
	)
}
