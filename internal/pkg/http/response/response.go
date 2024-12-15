package response

import (
	"log/slog"
	"net/http"
)

func RenderError(w http.ResponseWriter, err *HTTPError) {
	slog.Error(err.Error(), "error", err.Err)
	http.Error(w, err.Error(), err.Code)
}
