package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
)

type repo struct {
	cfg *config.DBConfig
	db  *sql.DB
}

func NewAuthRepo(cfg *config.DBConfig, conn *sql.DB) Authenticator {
	return &repo{
		cfg: cfg,
		db:  conn,
	}
}

func (r *repo) SignUp(ctx context.Context, params SignUpParams) (*User, error) {
	const q = "INSERT INTO users (email, password_hash, auth_method) VALUES ($1, $2, $3) RETURNING id, email, auth_method, created_at, updated_at"
	row := r.db.QueryRowContext(ctx, q, params.Email, params.Password, Basic)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.AuthMethod, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if db.IsUniqueViolation(err) {
			return nil, fmt.Errorf("user with email %s already exists: %w", params.Email, ErrEmailExists)
		}
		return nil, err
	}

	return &user, nil
}

func (r *repo) SignUpOAuth(ctx context.Context, provider string, id string) *sql.Row {
	const q = "INSERT INTO users (oauth_provider, oauth_id, auth_method) VALUES ($1, $2, $3) RETURNING id, email, auth_method, created_at, updated_at"
	return r.db.QueryRowContext(ctx, q, provider, id, OAuth)
}

func (r *repo) SignIn(ctx context.Context, email string) (string, error) {
	const q = "SELECT password_hash FROM users WHERE email = $1"
	row := r.db.QueryRowContext(ctx, q, email)

	var hash string
	if err := row.Scan(&hash); err != nil {
		return "", err
	}

	return hash, nil
}
