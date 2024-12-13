package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Model struct {
	ID        string          `json:"id,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime    `json:"deleted_at,omitempty"`
}
