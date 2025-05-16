package errors

const (
	InternalServerError = "INTERNAL_SERVER_ERROR"
)

type APIError struct {
	HttpCode int
	Type     string
	Message  string
}

type APIErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(httpCode int, errType string, errMsg string) *APIError {
	return &APIError{
		HttpCode: httpCode,
		Type:     errType,
		Message:  errMsg,
	}
}
