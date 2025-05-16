package errors

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/errors"
	"net/http"
)

const (
	InvalidCourtIDFormatError = "INVALID_COURT_ID_FORMAT_ERROR"
	CourtNotFoundError        = "COURT_NOT_FOUND_ERROR"
)

func NewInvalidCourtIDFormatError() *errors.APIError {
	return errors.NewAPIError(
		http.StatusBadRequest,
		InvalidCourtIDFormatError,
		"court id must be a valid uuid",
	)
}

func NewCourtNotFoundError(id uuid.UUID) *errors.APIError {
	return errors.NewAPIError(
		http.StatusNotFound,
		CourtNotFoundError,
		fmt.Sprintf("court with id: %s does not exist", id),
	)
}
