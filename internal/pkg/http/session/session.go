package session

import (
	"context"
	"net/http"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

// Session represents a single user session.
type Session interface {
	ID() []byte                             // Returns the unique session ID.
	Get(key string) any                     // Retrieves a value associated with the given key.
	Set(key string, value any) error        // Sets a value for the given key.
	Delete(key string) error                // Deletes a value associated with the given key.
	Flush() error                           // Persists session data to the storage.
	Destroy() error                         // Invalidates and removes the session.
	SetExpiry(duration time.Duration) error // Sets the session expiry time.
	Expiry() time.Time                      // Returns the session expiry time.
}

// Store represents the session storage backend.
type Store interface {
	Create(*http.Request) (Session, error)               // Creates a new session.
	Get(ctx context.Context, id []byte) (Session, error) // Retrieves a session by ID.
	Save(ctx context.Context, s Session) error           // Saves session data.
	Delete(ctx context.Context, id []byte) error         // Deletes a session by ID.
	Cleanup(ctx context.Context) error                   // Removes expired sessions.
}

// Manager manages sessions, providing high-level operations.
type Manager interface {
	Start(*http.Request) (Session, error)                // Starts a new session.
	Get(ctx context.Context, id []byte) (Session, error) // Retrieves a session by ID.
	Destroy(ctx context.Context, id []byte) error        // Destroys a session.
	UseStore(store Store)                                // Sets the session store.
	SetMaxLifeTime(duration time.Duration)               // Sets the maximum session lifetime.
	GC(ctx context.Context) error                        // Performs garbage collection of expired sessions.
}

// sessionData implements the Session interface.
type sessionData struct {
	id        []byte
	userID    string
	data      map[string]interface{}
	store     Store
	expiry    time.Time
	ipAddress string
	userAgent string
}

func (s *sessionData) ID() []byte {
	return s.id
}

func (s *sessionData) Get(key string) interface{} {
	return s.data[key]
}

func (s *sessionData) Set(key string, value interface{}) error {
	s.data[key] = value
	return nil
}

func (s *sessionData) Delete(key string) error {
	delete(s.data, key)
	return nil
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

func (s *sessionData) Expiry() time.Time {
	return s.expiry
}

type manager struct {
	store       Store
	maxLifeTime time.Duration
}

// NewManager creates a new session manager with the given options.
func NewManager(cfg config.SessionConfig) Manager {
	return &manager{maxLifeTime: cfg.SessionDuration}
}

func (m *manager) Start(r *http.Request) (Session, error) {
	sess, err := m.store.Create(r)
	if err != nil {
		return nil, err
	}
	err = sess.SetExpiry(m.maxLifeTime)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *manager) Get(ctx context.Context, id []byte) (Session, error) {
	return m.store.Get(ctx, id)
}

func (m *manager) Destroy(ctx context.Context, id []byte) error {
	return m.store.Delete(ctx, id)
}

func (m *manager) UseStore(store Store) {
	m.store = store
}

func (m *manager) SetMaxLifeTime(duration time.Duration) {
	m.maxLifeTime = duration
}

func (m *manager) GC(ctx context.Context) error {
	return m.store.Cleanup(ctx)
}
