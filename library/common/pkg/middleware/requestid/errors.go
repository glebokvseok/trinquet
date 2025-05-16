package reqidmw

import (
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"net/http"
)

const (
	InvalidRequestIdFormatError = "INVALID_REQUEST_ID_FORMAT_ERROR"
)

func NewInvalidRequestIdFormatError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidRequestIdFormatError,
		"Request id is in invalid format.",
	)
}
