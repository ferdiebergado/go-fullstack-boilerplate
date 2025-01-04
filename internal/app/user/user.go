package user

import (
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
)

type User struct {
	db.Model
	Email         string     `json:"email"`
	OAuthProvider *string    `json:"oauth_provider,omitempty"`
	OAuthID       *string    `json:"oauth_id,omitempty"`
	PasswordHash  *string    `json:"-"`
	AuthMethod    AuthMethod `json:"auth_method"`
}
