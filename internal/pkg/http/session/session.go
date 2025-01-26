package session

import (
	"context"
	"time"
)

// Session represents a single user session.
type Session interface {
	ID() []byte                             // Returns the unique session ID.
	UserID() *string                        // Returns the user ID.
	SetUserID(userID string) error          // Sets the user ID.
	Get(key string) any                     // Retrieves a value associated with the given key.
	Set(key string, value any) error        // Sets a value for the given key.
	Delete(key string) error                // Deletes a value associated with the given key.
	Flash(key string) any                   // Retrieves a value and then deletes it from the store.
	Flush() error                           // Persists session data to the storage.
	Destroy() error                         // Invalidates and removes the session.
	SetExpiry(duration time.Duration) error // Sets the session expiry time.
	Expiry() time.Time                      // Returns the session expiry time.
}

// sessionData implements the Session interface.
type sessionData struct {
	id        []byte
	userID    *string
	data      map[string]any
	store     Store
	expiry    time.Time
	ipAddress string
	userAgent string
}

var _ Session = (*sessionData)(nil)

func (s *sessionData) ID() []byte {
	return s.id
}

func (s *sessionData) UserID() *string {
	return s.userID
}

func (s *sessionData) Get(key string) any {
	return s.data[key]
}

func (s *sessionData) Set(key string, value any) error {
	s.data[key] = value
	return nil
}

func (s *sessionData) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *sessionData) Flash(key string) any {
	val := s.data[key]
	delete(s.data, key)
	return val
}

func (s *sessionData) Flush() error {
	return s.store.Save(context.Background(), s)
}

func (s *sessionData) Destroy() error {
	return s.store.Delete(context.Background(), s.ID())
}

func (s *sessionData) SetExpiry(duration time.Duration) error {
	s.expiry = time.Now().Add(duration)
	return s.Flush()
}

func (s *sessionData) SetUserID(userID string) error {
	s.userID = &userID
	return s.Flush()
}

func (s *sessionData) Expiry() time.Time {
	return s.expiry
}
