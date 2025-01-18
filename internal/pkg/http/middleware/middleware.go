package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
)

type Middleware func(http.Handler) http.Handler

func SessionMiddleware(cfg config.HTTPServerConfig, sessMgr session.Manager) Middleware {
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID")

		type redirectData struct {
			RedirectPath string `json:"redirect_path"`
		}

		if userID == nil {
			data := &response.APIResponse[redirectData]{
				Message: http.StatusText(http.StatusUnauthorized),
				Data: &redirectData{
					RedirectPath: "/signin",
				},
			}
			response.RenderJSON(w, http.StatusFound, data)
			return
		}

		next.ServeHTTP(w, r)
	})
}
