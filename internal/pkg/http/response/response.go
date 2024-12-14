package response

import (
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
)

func RenderError(w http.ResponseWriter, err *HTTPError, logger *logging.Logger) {
	logger.Error(err.Error(), slog.String("error", err.Err.Error()))
	http.Error(w, err.Error(), err.Code)
}
