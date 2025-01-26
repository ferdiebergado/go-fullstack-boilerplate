package session

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/request"
)

// Manager manages sessions, providing high-level operations.
type Manager interface {
	Start(r *http.Request) (Session, error)       // Starts a new session.
	Get(r *http.Request) (Session, error)         // Retrieves a session by ID.
	Destroy(ctx context.Context, id []byte) error // Destroys a session.
	UseStore(store Store)                         // Sets the session store.
	SetMaxLifeTime(duration time.Duration)        // Sets the maximum session lifetime.
	GC(ctx context.Context) error                 // Performs garbage collection of expired sessions.
	Config() *config.SessionConfig
}

type manager struct {
	store       Store
	maxLifeTime time.Duration
	config      *config.SessionConfig
}

// NewManager creates a new session manager with the given options.
func NewManager(cfg config.SessionConfig) Manager {
	return &manager{maxLifeTime: cfg.SessionDuration, config: &cfg}
}

func (m *manager) Start(r *http.Request) (Session, error) {
	slog.Debug("request.getipaddress", "val", request.GetIPAddress(r))
	sess, err := m.store.Create(r.Context(), request.GetIPAddress(r), r.UserAgent())
	if err != nil {
		return nil, err
	}

	if err := sess.SetExpiry(m.maxLifeTime); err != nil {
		return nil, err
	}

	return sess, nil
}

func (m *manager) Get(r *http.Request) (Session, error) {
	sessCookie, err := r.Cookie(m.config.SessionName)
	if err != nil {
		return nil, err
	}

	sessID, err := base64.RawURLEncoding.DecodeString(sessCookie.Value)
	if err != nil {
		return nil, err
	}

	return m.store.Get(r.Context(), sessID)
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

func (m *manager) Config() *config.SessionConfig {
	return m.config
}
