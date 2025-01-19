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
			sessionData, err := sessMgr.Fetch(r)

			if err != nil {
				slog.Debug("No data for session")
				next.ServeHTTP(w, r)
				return
			}

			slog.Debug("Session", "data", sessionData, "user_id", sessionData.UserID)

			ctx := WithUser(r.Context(), sessionData.UserID)
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

				sessionData, err := sessMgr.Fetch(r)

				if err != nil {
					slog.Error("no session data")
				} else {
					sessionData.Flash = map[string]string{
						"intendedUrl": r.URL.Path,
					}

					sessionID, err := sessMgr.SessionID(r)

					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					sessMgr.Save(r.Context(), sessionID, *sessionData)
				}

				http.Redirect(w, r, redirectPath, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
