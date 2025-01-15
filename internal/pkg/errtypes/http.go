package errtypes

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/validation"
)

type HTTPError struct {
	Msg  string `json:"msg"`
	Err  error  `json:"err"`
	Code int    `json:"code"`
}

func (e HTTPError) Error() string {
	return e.Msg
}

func ServerError(err error) *HTTPError {
	return &HTTPError{
		Msg:  "Something went wrong.",
		Err:  err,
		Code: http.StatusInternalServerError,
	}
}

func ServerUnavailableError(err error) *HTTPError {
	return &HTTPError{
		Msg:  "Something went wrong.",
		Err:  err,
		Code: http.StatusServiceUnavailable,
	}
}

func BadRequest(err error) *HTTPError {
	return &HTTPError{
		Msg:  "Cannot read the request.",
		Err:  err,
		Code: http.StatusBadRequest,
	}
}

func ValidationError(inputErr validation.Error) *HTTPError {
	return &HTTPError{
		Msg:  inputErr.Error(),
		Err:  &inputErr,
		Code: http.StatusUnprocessableEntity,
	}
}

func AuthenticationError(err error) *HTTPError {
	return &HTTPError{
		Msg:  err.Error(),
		Err:  err,
		Code: http.StatusUnauthorized,
	}
}

func JSONEncodeError(err error) *HTTPError {
	return &HTTPError{
		Msg:  "Failed to encode json.",
		Err:  err,
		Code: http.StatusInternalServerError,
	}
}
