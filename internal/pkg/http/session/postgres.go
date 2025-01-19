package session

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/security"
)

type DatabaseSession struct {
	cfg   config.SessionConfig
	store *sql.DB
}

func NewDatabaseSession(cfg config.SessionConfig, db *sql.DB) Manager {
	gob.Register(Data{})

	return &DatabaseSession{
		cfg:   cfg,
		store: db,
	}
}

func (d *DatabaseSession) Save(ctx context.Context, sessionID string, sessionData Data) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(sessionData); err != nil {
		return fmt.Errorf("encode session data: %w", err)
	}

	expiryTime := time.Now().Add(d.cfg.SessionDuration)

	_, err := d.store.ExecContext(ctx,
		`INSERT INTO user_sessions (session_id, session_data, last_activity, expiry_time)
                VALUES ($1, $2, NOW(), $3)
                ON CONFLICT (session_id) DO UPDATE SET session_data = $2, last_activity = NOW(), expiry_time = $3`,
		sessionID, buf.Bytes(), expiryTime)

	if err != nil {
		return fmt.Errorf("save session data: %w", err)
	}

	return nil
}

func (d *DatabaseSession) Fetch(r *http.Request) (*Data, error) {
	sessionID, err := d.SessionID(r)

	if err != nil {
		return nil, err
	}

	var sessionData []byte
	err = d.store.QueryRowContext(r.Context(), "SELECT session_data FROM user_sessions WHERE session_id = $1", sessionID).Scan(&sessionData)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get session data: %w", err)
	}

	var data Data
	err = gob.NewDecoder(bytes.NewReader(sessionData)).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode session data: %w", err)
	}

	return &data, nil
}

func (d *DatabaseSession) Destroy(r *http.Request) error {
	sessionID, err := d.SessionID(r)

	if err != nil {
		return err
	}

	_, err = d.store.ExecContext(r.Context(), "DELETE FROM user_sessions WHERE session_id = $1", sessionID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (d *DatabaseSession) SessionID(r *http.Request) (string, error) {
	session, err := r.Cookie(d.cfg.SessionName)
	var sessionID string
	if err != nil {
		sessionID, err = security.GenerateRandomBytesEncoded(64)

		if err != nil {
			return "", fmt.Errorf("generate session id: %w", err)
		}
	} else {
		sessionID = session.Value
	}

	return sessionID, nil
}
