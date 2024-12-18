package response

import (
	"errors"
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

func TemplateNotFoundError(name string) *HTTPError {
	msg := "template does not exists: " + name

	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Err:     errors.New(msg),
	}
}

func (e *HTTPError) Error() string {
	return e.Message
}
