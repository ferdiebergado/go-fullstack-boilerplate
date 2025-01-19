package auth

import (
	"context"
	"errors"
)

type ctxKey string

const userKey ctxKey = "userID"

var ErrUserNotInContext = errors.New("no user in context")

func WithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

func FromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userKey).(string)
	if !ok {
		return "", ErrUserNotInContext
	}

	return userID, nil
}
