package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/auth"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

type authenticator struct {
	SignUpFn func(context.Context, auth.SignUpParams) (*user.User, error)
	SignInFn func(context.Context, string) (*auth.SignInResult, error)
}

func (r *authenticator) SignUp(ctx context.Context, params auth.SignUpParams) (*user.User, error) {
	if r.SignUpFn != nil {
		return r.SignUpFn(ctx, params)
	}

	return nil, nil
}

func (r *authenticator) SignIn(ctx context.Context, email string) (*auth.SignInResult, error) {
	if r.SignInFn != nil {
		return r.SignInFn(ctx, email)
	}

	return nil, nil
}

func TestServiceSignUpHappyPath(t *testing.T) {
	repo := &authenticator{
		SignUpFn: func(_ context.Context, _ auth.SignUpParams) (*user.User, error) {
			return &user.User{
				Email:      testEmail,
				AuthMethod: user.BasicAuth,
			}, nil
		},
	}

	service := auth.NewAuthService(&config.Config{}, repo)

	got, err := service.SignUp(context.Background(), signupParams)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Email != signupParams.Email {
		t.Errorf("want: %v but got: %v", signupParams.Email, got.Email)
	}

	if string(got.AuthMethod) != string(user.BasicAuth) {
		t.Errorf("want: %v but got: %v", user.BasicAuth, got.AuthMethod)
	}
}

func TestServiceSignUpDuplicateEmail(t *testing.T) {
	repo := &authenticator{
		SignUpFn: func(_ context.Context, _ auth.SignUpParams) (*user.User, error) {
			return nil, &auth.EmailExistsError{Email: testEmail}
		},
	}

	service := auth.NewAuthService(&config.Config{}, repo)

	_, err := service.SignUp(context.Background(), signupParams)

	if err == nil {
		t.Error("expected an error but got nil")
	}

	var emailExistsErr *auth.EmailExistsError
	if !errors.As(err, &emailExistsErr) {
		t.Errorf("want: %v but got: %v", emailExistsErr, err)
	}
}

func TestServiceSignUpRepoError(t *testing.T) {
	repo := &authenticator{
		SignUpFn: func(_ context.Context, _ auth.SignUpParams) (*user.User, error) {
			return nil, errors.New("repo error")
		},
	}

	service := auth.NewAuthService(&config.Config{}, repo)

	_, err := service.SignUp(context.Background(), signupParams)

	if err == nil {
		t.Error("expected an error but got nil")
	}
}
