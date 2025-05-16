package errors

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

const (
	RequestBodyParsingError = "REQUEST_BODY_PARSING_ERROR"

	InvalidRequestBodyFormatError = "INVALID_REQUEST_BODY_FORMAT_ERROR"
)

func NewRequestBodyParsingError(err error) *APIError {
	errMsg := err.Error()
	if httpErr := (*echo.HTTPError)(nil); errors.As(err, &httpErr) {
		errMsg = fmt.Sprintf("%+v", httpErr.Message)
	}

	return NewAPIError(
		http.StatusBadRequest,
		RequestBodyParsingError,
		fmt.Sprintf("following error occured while parsing request body: %s", errMsg),
	)
}

func NewInvalidRequestBodyFormatError(err error) *APIError {
	return NewAPIError(
		http.StatusBadRequest,
		InvalidRequestBodyFormatError,
		fmt.Sprintf("following errors occured while validating request body:\n%s", err.Error()),
	)
}
