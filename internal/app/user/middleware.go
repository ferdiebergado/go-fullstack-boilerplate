package user

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
)

const redirectPath = "/signin"

func AuthMiddleware(next http.Handler) http.Handler {
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

			http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
