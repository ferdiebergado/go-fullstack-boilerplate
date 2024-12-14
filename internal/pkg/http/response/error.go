package response

import (
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

const ServerErrorMessage = "An error occurred."

func ServerError(err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: ServerErrorMessage,
		Err:     err,
	}
}

func (e *HTTPError) Error() string {
	return e.Message
}
