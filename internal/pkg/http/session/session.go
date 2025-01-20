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
	StoreSession(context.Context, string, Data) error

	// Loads the session data from the request.
	LoadSession(*http.Request) (*Data, error)

	// Extracts the session id from the request.
	ExtractSessionID(*http.Request) (string, error)

	// Deletes the session from the request.
	DestroySession(*http.Request) error
}
