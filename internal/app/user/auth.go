package user

import "context"

type AuthMethod string

const (
	Basic AuthMethod = "email/password"
	OAuth AuthMethod = "oauth"
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

type Authenticator interface {
	SignUp(context.Context, SignUpParams) (*User, error)
	SignIn(context.Context, string) (string, error)
}
