package errtypes

import "net/http"

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
