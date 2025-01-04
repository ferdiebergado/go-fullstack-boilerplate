package user

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
)

type AuthMethod string

const (
	Basic AuthMethod = "email/password"
	OAuth AuthMethod = "oauth"
)

type AuthData struct {
	Email                string     `json:"email,omitempty"`
	Password             string     `json:"password,omitempty"`
	PasswordConfirmation string     `json:"password_confirmation,omitempty"`
	OAuthProvider        string     `json:"o_auth_provider,omitempty"`
	OAuthID              string     `json:"o_auth_id,omitempty"`
	AuthMethod           AuthMethod `json:"auth_method"`
}

type BasicAuthParams struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type Authenticator interface {
	SignUp(context.Context, AuthData) (User, error)
	SignInBasic(context.Context, string) (*SignInParams, error)
}

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

func (r *repo) SignUp(ctx context.Context, data AuthData) (User, error) {
	var row *sql.Row

	switch data.AuthMethod {
	case Basic:
		row = r.signUpBasic(ctx, data.Email, data.Password)
	case OAuth:
		row = r.signUpOAuth(ctx, data.OAuthProvider, data.OAuthID)
	}

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.OAuthProvider, &user.AuthMethod, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if db.IsUniqueViolation(err) {
			return user, fmt.Errorf("user with email %s already exists: %w", data.Email, ErrEmailExists)
		}
		return user, err
	}

	return user, nil
}

const signUpquery = "INSERT INTO users (%s, %s, %s) VALUES ($1, $2, $3) RETURNING id, email, oauth_provider, auth_method, created_at, updated_at"

func (r *repo) signUpBasic(ctx context.Context, email string, passwordHash string) *sql.Row {
	q := fmt.Sprintf(signUpquery, "email", "password_hash", "auth_method")
	slog.Debug("signupbasic q", "q", q, "email", email, "password", passwordHash)
	return r.db.QueryRowContext(ctx, q, email, passwordHash, Basic)
}

func (r *repo) signUpOAuth(ctx context.Context, provider string, id string) *sql.Row {
	q := fmt.Sprintf(signUpquery, "oauth_provider", "oauth_id", "auth_method")
	return r.db.QueryRowContext(ctx, q, provider, id, OAuth)
}

type SignInParams struct {
	Email    string
	Password string
}

func (r *repo) SignInBasic(ctx context.Context, email string) (*SignInParams, error) {
	const signInquery = "SELECT email, password_hash FROM users WHERE email = $1"
	row := r.db.QueryRowContext(ctx, signInquery, email)

	var params SignInParams
	if err := row.Scan(&params.Email, &params.Password); err != nil {
		return nil, err
	}

	return &params, nil
}
