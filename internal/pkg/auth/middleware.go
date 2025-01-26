package auth

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/middleware"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/request"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/security"
)

const redirectPath = "/signin"

func SessionMiddleware(sessMgr session.Manager) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("request headers", "x_real_ip", r.Header.Get("X-Real-IP"), "x_forwarded_for", r.Header.Get("X-Forwarded-For"), "remote_address", r.RemoteAddr, "request_getipaddress", request.GetIPAddress(r))
			sessData, err := sessMgr.Get(r)

			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				slog.Error("Unable to retrieve session", "reason", err)
				next.ServeHTTP(w, r)
				return
			}

			if sessData == nil {
				slog.Debug("No session cookie. Starting a new session.")
				sessData, err := sessMgr.Start(r)

				if err != nil {
					slog.Error("Unable to start the session.", "reason", err)
					next.ServeHTTP(w, r)
					return
				}

				cfg := sessMgr.Config()
				sessID := base64.RawURLEncoding.EncodeToString(sessData.ID())

				http.SetCookie(w, &http.Cookie{
					Name:     cfg.SessionName,
					Value:    sessID,
					Expires:  time.Now().Add(cfg.SessionDuration),
					HttpOnly: true,
					SameSite: cfg.SameSite,
					Path:     "/",
				})

				slog.Debug("Session cookie set.", "session_id", sessID)

				next.ServeHTTP(w, r)
				return
			}

			userID := sessData.UserID()

			if userID != nil {
				slog.Debug("Session", "user_id", *userID)
				ctx := WithUser(r.Context(), *userID)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireUserMiddleware(sessMgr session.Manager) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := FromContext(r.Context())

			type redirectData struct {
				RedirectPath string `json:"redirect_path"`
			}

			if userID == nil {
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

				sessionData, err := sessMgr.Get(r)

				if err != nil {
					http.Redirect(w, r, redirectPath, http.StatusSeeOther)
					return
				}

				if sessionData != nil {
					if err := sessionData.Set("intendedUrl", r.URL.Path); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if err := sessionData.Flush(); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}

				http.Redirect(w, r, redirectPath, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CSRFMiddleware(cfg config.SessionConfig) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := r.Cookie(cfg.CSRFName)

			if err != nil {
				csrf, err := security.GenerateRandomBytesEncoded(64)

				if err != nil {
					serverError := errtypes.ServerError(err)
					response.RenderError(w, r, serverError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     cfg.CSRFName,
					Value:    csrf,
					Expires:  time.Now().Add(cfg.SessionDuration),
					HttpOnly: false,
					SameSite: cfg.SameSite,
					Path:     "/",
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}
