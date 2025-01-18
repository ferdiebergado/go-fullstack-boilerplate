package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/security"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/validation"
)

type service struct {
	authenticator Authenticator
	cfg           *config.Config
}

type AuthService interface {
	SignUp(context.Context, SignUpParams) (*User, error)
	SignIn(context.Context, SignInParams) (string, error)
}

func NewAuthService(cfg *config.Config, authenticator Authenticator) AuthService {
	return &service{
		authenticator: authenticator,
		cfg:           cfg,
	}
}

var ErrEmailExists = errors.New("duplicate email")
var ErrUserPassInvalid = errors.New("invalid username or password")

// Signs up a user using email and password
func (s *service) SignUp(ctx context.Context, params SignUpParams) (*User, error) {
	form := validation.NewForm(params)
	form.Required("Email", "Password", "PasswordConfirmation")
	form.PasswordsMatch("Password", "PasswordConfirmation")
	form.IsEmail("Email")

	if !form.IsValid() {
		return nil, &validation.Error{
			Errors: form.Error.Errors,
		}
	}

	hash, err := security.GenerateHash(params.Password)

	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	params.Password = hash

	return s.authenticator.SignUp(ctx, params)
}

// Signs in a user using email and password
func (s *service) SignIn(ctx context.Context, params SignInParams) (string, error) {
	form := validation.NewForm(params)
	form.Required("Email", "Password")
	form.IsEmail("Email")

	if !form.IsValid() {
		return "", &validation.Error{
			Errors: form.Error.Errors,
		}
	}

	result, err := s.authenticator.SignIn(ctx, params.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("find user: %w", ErrUserPassInvalid)
		}

		return "", err
	}

	slog.Debug("sign in", "hash", result.Hash)

	match, err := security.VerifyPassword(params.Password, result.Hash)

	if err != nil {
		return "", err
	}

	if !match {
		return "", fmt.Errorf("verify password: %w", ErrUserPassInvalid)
	}

	return result.ID, nil
}
