package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/validation"
	gkitResponse "github.com/ferdiebergado/gopherkit/http/response"
)

type PageData struct {
	Title    string
	Subtitle string
}

type APIResponse[T any] struct {
	Message string            `json:"message,omitempty"`
	Errors  validation.Errors `json:"errors,omitempty"`
	Data    *T                `json:"data,omitempty"`
}

func RenderError(w http.ResponseWriter, r *http.Request, err *errtypes.HTTPError) {
	slog.Error(err.Msg, "error", err.Err)

	if r.Header.Get("content-type") != "application/json" {
		http.Error(w, err.Error(), err.Code)
		return
	}

	res := &APIResponse[any]{
		Message: err.Error(),
	}

	var valErr *validation.Error
	if errors.As(err.Err, &valErr) {
		res.Errors = valErr.Errors
	}

	RenderJSON(w, err.Code, res)
}

func RenderJSON[T any](w http.ResponseWriter, status int, data *T) {
	if err := gkitResponse.JSON(w, status, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
