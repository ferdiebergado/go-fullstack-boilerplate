package errtypes

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/validation"
)

type HTTPError struct {
	*AppError
	Code int
}

func ServerError(err error) *HTTPError {
	return &HTTPError{
		AppError: &AppError{
			Description: "An error occurred.",
			Err:         err,
			Severity:    Critical,
		},
		Code: http.StatusInternalServerError,
	}
}

func ValidationError(inputErr validation.Error) *HTTPError {
	return &HTTPError{
		AppError: &AppError{
			Description: inputErr.Error(),
			Err:         &inputErr,
			Severity:    Low,
		},
		Code: http.StatusUnprocessableEntity,
	}
}

func JSONEncodeError(err error) *HTTPError {
	return &HTTPError{
		AppError: &AppError{Description: "failed to encode json", Err: err, Severity: High},
		Code:     http.StatusInternalServerError,
	}
}
