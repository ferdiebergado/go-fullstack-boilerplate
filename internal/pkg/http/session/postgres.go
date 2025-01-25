package session

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/request"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/security"
)

// postgresStore implements the Store interface using PostgreSQL.
type postgresStore struct {
	cfg config.SessionConfig
	db  *sql.DB
}

// NewPostgresStore creates a new PostgresStore.
func NewPostgresStore(cfg config.SessionConfig, db *sql.DB) Store {
	return &postgresStore{
		cfg: cfg,
		db:  db,
	}
}

func (s *postgresStore) Create(r *http.Request) (Session, error) {
	id, err := security.GenerateRandomBytes(64)

	if err != nil {
		return nil, err
	}

	sess := &sessionData{
		id:        id,
		data:      make(map[string]interface{}),
		store:     s,
		expiry:    time.Now().Add(s.cfg.SessionDuration),
		ipAddress: request.GetIPAddress(r),
		userAgent: r.UserAgent(),
	}

	if err := s.Save(r.Context(), sess); err != nil {
		return nil, fmt.Errorf("failed to save new session: %w", err)
	}

	return sess, nil
}

func (s *postgresStore) Get(ctx context.Context, id []byte) (Session, error) {
	var data []byte
	var expiry time.Time
	err := s.db.QueryRowContext(ctx, "SELECT data, expiry FROM sessions WHERE id = $1", id).Scan(&data, &expiry)
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

func (s *postgresStore) Save(ctx context.Context, sess Session) error {
	sessData, ok := sess.(*sessionData)

	if !ok {
		return fmt.Errorf("failed to type assert session to sessionData: %v", sess)
	}

	data, err := json.Marshal(sessData.data)
	if err != nil {
		return fmt.Errorf("failed to encode session data: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO sessions (id, user_id, ip_address, user_agent, data, expiry) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE SET data = $5, expiry = $6, last_activity = NOW()", sess.ID(), sessData.userID, sessData.ipAddress, sessData.userAgent, data, sess.Expiry())

	return err
}

func (s *postgresStore) Delete(ctx context.Context, id []byte) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}

func (s *postgresStore) Cleanup(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE expiry < NOW()")
	return err
}
