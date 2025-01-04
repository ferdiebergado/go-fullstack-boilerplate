package db

import (
	"errors"
	"strings"
)

var ErrRowClose = errors.New("failed to close the rows result set")
var ErrRowScan = errors.New("error occurred while scanning the row into the destination variables")
var ErrRowIteration = errors.New("error encountered during row iteration, possibly due to a database or connection issue")
var ErrModelNotFound = errors.New("model not found")

// isUniqueViolation checks if an error is a unique constraint violation based on SQLSTATE code 23505.
func IsUniqueViolation(err error) bool {
	// PostgreSQL drivers usually embed the SQLSTATE in the error string
	return strings.Contains(err.Error(), "23505")
}
