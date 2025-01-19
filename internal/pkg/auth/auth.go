package auth

import (
	"context"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
)

type SignUpParams struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type SignInParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OAuthParams struct {
	OAuthProvider string `json:"oauth_provider"`
	OAuthID       string `json:"oauth_id"`
}

type SignInResult struct {
	ID   string
	Hash string
}

type Authenticator interface {
	SignUp(context.Context, SignUpParams) (*user.User, error)
	SignIn(context.Context, string) (*SignInResult, error)
}
