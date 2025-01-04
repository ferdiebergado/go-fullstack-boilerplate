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
	SignUp(context.Context, AuthData) (User, error)
	SignIn(context.Context, AuthData) error
}

func NewAuthService(cfg *config.Config, authenticator Authenticator) AuthService {
	return &service{
		authenticator: authenticator,
		cfg:           cfg,
	}
}

var ErrEmailExists = errors.New("duplicate email")
var ErrUserPassInvalid = errors.New("invalid username or password")

func (s *service) SignUp(ctx context.Context, data AuthData) (User, error) {
	form := validation.NewForm(data)
	form.Required("Email", "Password", "PasswordConfirmation")
	form.PasswordsMatch("Password", "PasswordConfirmation")
	form.IsEmail("Email")

	if !form.IsValid() {
		return User{}, &validation.InputError{
			Errors: form.Errors,
		}
	}

	if data.AuthMethod == Basic {
		hash, err := security.GenerateHash(data.Password)

		if err != nil {
			return User{}, fmt.Errorf("hash password: %w", err)
		}

		data.Password = hash
	}

	return s.authenticator.SignUp(ctx, data)
}

func (s *service) SignIn(ctx context.Context, data AuthData) error {
	params, err := s.authenticator.SignInBasic(ctx, data.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("find user: %w", ErrUserPassInvalid)
		}

		return err
	}

	slog.Debug("params", "password", params.Password)

	match, err := security.VerifyPassword(data.Password, params.Password)

	if err != nil {
		return err
	}

	if !match {
		return fmt.Errorf("verify password: %w", ErrUserPassInvalid)
	}

	return nil
}
