package session

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/security"
)

// Store represents the session storage backend.
type Store interface {
	Create(ctx context.Context, ipAddress string, userAgent string) (Session, error) // Creates a new session.
	Get(ctx context.Context, id []byte) (Session, error)                             // Retrieves a session by ID.
	Save(ctx context.Context, s Session) error                                       // Saves session data.
	Delete(ctx context.Context, id []byte) error                                     // Deletes a session by ID.
	Cleanup(ctx context.Context) error                                               // Removes expired sessions.
}

// dbStore implements the Store interface using PostgreSQL.
type dbStore struct {
	cfg config.SessionConfig
	db  *sql.DB
}

// NewPostgresStore creates a new PostgresStore.
func NewPostgresStore(cfg config.SessionConfig, db *sql.DB) Store {
	return &dbStore{
		cfg: cfg,
		db:  db,
	}
}

func (s *dbStore) Create(ctx context.Context, ipAddress string, userAgent string) (Session, error) {
	id, err := security.GenerateRandomBytes(64)

	if err != nil {
		return nil, err
	}

	sess := &sessionData{
		id:        id,
		userID:    nil,
		data:      make(map[string]any),
		store:     s,
		expiry:    time.Now().Add(s.cfg.SessionDuration),
		ipAddress: ipAddress,
		userAgent: userAgent,
	}

	if err := s.Save(ctx, sess); err != nil {
		return nil, fmt.Errorf("failed to save new session: %w", err)
	}

	return sess, nil
}

func (s *dbStore) Get(ctx context.Context, id []byte) (Session, error) {
	var data []byte
	var expiry time.Time
	var userID *string
	err := s.db.QueryRowContext(ctx, "SELECT user_id, data, expiry_time FROM sessions WHERE id = $1", id).Scan(&userID, &data, &expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Session not found
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if expiry.Before(time.Now()) {
		if err := s.Delete(ctx, id); err != nil {
			return nil, err
		}
		//cleanup expired session
		return nil, nil
	}

	sess := &sessionData{
		id:     id,
		userID: userID,
		store:  s,
		expiry: expiry,
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &sess.data); err != nil {
			return nil, fmt.Errorf("failed to decode session data: %w", err)
		}
	} else {
		sess.data = make(map[string]interface{})
	}

	return sess, nil
}

const saveQuery = `
INSERT INTO sessions (id, user_id, ip_address, user_agent, data, expiry_time)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE
SET data = $5, expiry_time = $6, last_activity = NOW()
`

func (s *dbStore) Save(ctx context.Context, sess Session) error {
	sessData, ok := sess.(*sessionData)

	if !ok {
		return fmt.Errorf("failed to type assert session to sessionData: %v", sess)
	}

	data, err := json.Marshal(sessData.data)
	if err != nil {
		return fmt.Errorf("failed to encode session data: %w", err)
	}

	_, err = s.db.ExecContext(ctx, saveQuery, sess.ID(), sessData.userID, sessData.ipAddress, sessData.userAgent, data, sess.Expiry())

	return err
}

func (s *dbStore) Delete(ctx context.Context, id []byte) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}

func (s *dbStore) Cleanup(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE expiry_time < NOW()")
	return err
}
