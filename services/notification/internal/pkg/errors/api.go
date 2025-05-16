package errors

import (
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"net/http"
)

const (
	InvalidTimestampFormatError = "INVALID_TIMESTAMP_FORMAT_ERROR"
)

func NewInvalidTimestampFormatError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidTimestampFormatError,
		"cursor must be an int64",
	)
}
