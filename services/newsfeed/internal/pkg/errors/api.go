package errors

import (
	"fmt"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
	"net/http"
)

const (
	InvalidMediaTypeFormatError = "INVALID_MEDIA_TYPE_FORMAT_ERROR"
	UnsupportedMediaTypeError   = "UNSUPPORTED_MEDIA_TYPE_ERROR"

	InvalidPostIDFormatError = "INVALID_POST_ID_FORMAT_ERROR"

	InvalidCursorFormatError = "INVALID_CURSOR_POSITION_FORMAT_ERROR"
)

func NewInvalidMediaTypeError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidMediaTypeFormatError,
		"media type must be a string",
	)
}

func NewUnsupportedMediaTypeError(mediaType models.MediaType) *errors.APIError {
	return errors.NewAPIError(
		http.StatusNotFound,
		UnsupportedMediaTypeError,
		fmt.Sprintf("unsuppored media type: %s", mediaType),
	)
}

func NewInvalidPostIDFormatError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidPostIDFormatError,
		"post id must be a valid uuid",
	)
}

func NewInvalidCursorFormatError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidCursorFormatError,
		"cursor must be an int64",
	)
}
