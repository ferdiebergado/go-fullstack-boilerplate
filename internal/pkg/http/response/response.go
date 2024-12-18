package response

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
)

func RenderError(w http.ResponseWriter, err *errtypes.HTTPError) {
	slog.Error(err.Error(), "error", err.Err, "severity", err.Severity)
	http.Error(w, err.Error(), err.Code)
}
