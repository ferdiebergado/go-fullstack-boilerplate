package response

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	gkitResponse "github.com/ferdiebergado/gopherkit/http/response"
)

type PageData struct {
	Title    string
	Subtitle string
}

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

func RenderError[T any](w http.ResponseWriter, err *errtypes.HTTPError, data *T) {
	slog.Error(err.Error(), "error", err.Err, "severity", err.Severity)

	if data != nil {
		if err := gkitResponse.JSON(w, err.Code, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	http.Error(w, err.Error(), err.Code)
}

func RenderJSON[T any](w http.ResponseWriter, status int, data *T) {
	if err := gkitResponse.JSON(w, status, data); err != nil {
		RenderError[T](w, errtypes.JSONEncodeError(err), nil)
	}
}
