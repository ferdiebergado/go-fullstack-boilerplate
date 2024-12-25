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

func JSONEncodeError(err error) *HTTPError {
	return &HTTPError{
		AppError: &AppError{Description: "failed to encode json", Err: err, Severity: High},
		Code:     http.StatusInternalServerError,
	}
}
