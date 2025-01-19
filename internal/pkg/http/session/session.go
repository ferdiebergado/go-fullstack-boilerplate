package session

import (
	"context"
	"net/http"
)

type Data struct {
	UserID string
	Flash  map[string]string
}

type Manager interface {
	// Saves a session.
	Save(context.Context, string, Data) error

	// Retrieves the session data from the request.
	Fetch(*http.Request) (*Data, error)

	// Retrieves the session id from the request.
	SessionID(*http.Request) (string, error)

	// Deletes the session data from the request.
	Destroy(*http.Request) error
}
