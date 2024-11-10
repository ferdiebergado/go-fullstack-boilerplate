package db

import (
	"encoding/json"
	"time"
)

type Model struct {
	ID        int64           `json:"id,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
}
