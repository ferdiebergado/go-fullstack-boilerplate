package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
)

type Middleware func(http.Handler) http.Handler

func SessionMiddleware(cfg config.SessionConfig, sessMgr session.Manager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie(cfg.SessionName)

			if err != nil {
				slog.Info("No session cookie")
				next.ServeHTTP(w, r)
				return
			}

			userID, err := sessMgr.Session(session.Value)

			if err != nil {
				slog.Info("No user from session")
				next.ServeHTTP(w, r)
				return
			}

			slog.Info("Session", "session", session, "user_id", userID)

			ctx := user.WithUser(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
