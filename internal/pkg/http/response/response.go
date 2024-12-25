package response

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	gkitResponse "github.com/ferdiebergado/gopherkit/http/response"
)

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
