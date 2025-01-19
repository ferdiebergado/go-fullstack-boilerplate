package auth

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/middleware"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
)

const redirectPath = "/signin"

func SessionMiddleware(cfg config.SessionConfig, sessMgr session.Manager) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie(cfg.SessionName)

			if err != nil {
				slog.Debug("No session cookie")
				next.ServeHTTP(w, r)
				return
			}

			userID, err := sessMgr.Session(session.Value)

			if err != nil {
				slog.Debug("No user from session")
				next.ServeHTTP(w, r)
				return
			}

			slog.Debug("Session", "session", session, "user_id", userID)

			ctx := WithUser(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireUserMiddleware(sessMgr session.Manager) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := FromContext(r.Context())

			type redirectData struct {
				RedirectPath string `json:"redirect_path"`
			}

			if userID == "" || err != nil {
				if r.Header.Get("content-type") == "application/json" {
					data := &response.APIResponse[redirectData]{
						Message: "Login is required to access this resource.",
						Data: &redirectData{
							RedirectPath: redirectPath,
						},
					}
					response.RenderJSON(w, http.StatusFound, data)
					return
				}

				if err := sessMgr.Save("intended_url", r.URL.Path); err != nil {
					slog.Info("cannot save url to the session")
					next.ServeHTTP(w, r)
					return
				}

				http.Redirect(w, r, redirectPath, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
