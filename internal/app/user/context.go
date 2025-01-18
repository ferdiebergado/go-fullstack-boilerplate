package user

import "context"

type ctxKey string

const userKey ctxKey = "userID"

func WithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey, userID)
}
