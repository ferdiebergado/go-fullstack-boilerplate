package auth

import (
	"context"
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
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

const signUpQuery = `
INSERT INTO users (email, password_hash, auth_method)
VALUES ($1, $2, $3)
RETURNING id, email, auth_method, created_at, updated_at
`

func (r *repo) SignUp(ctx context.Context, params SignUpParams) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, signUpQuery, params.Email, params.Password, user.Basic)

	var user user.User
	if err := row.Scan(&user.ID, &user.Email, &user.AuthMethod, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if db.IsUniqueViolation(err) {
			return nil, &EmailExistsError{Email: params.Email}
		}
		return nil, err
	}

	return &user, nil
}

func (r *repo) SignUpOAuth(ctx context.Context, provider string, id string) *sql.Row {
	const q = "INSERT INTO users (oauth_provider, oauth_id, auth_method) VALUES ($1, $2, $3) RETURNING id, email, auth_method, created_at, updated_at"
	return r.db.QueryRowContext(ctx, q, provider, id, user.OAuth)
}

const singInQuery = `
SELECT id, password_hash FROM users
WHERE email = $1
`

func (r *repo) SignIn(ctx context.Context, email string) (*SignInResult, error) {
	row := r.db.QueryRowContext(ctx, singInQuery, email)

	var result SignInResult
	if err := row.Scan(&result.ID, &result.Hash); err != nil {
		return nil, err
	}

	return &result, nil
}
