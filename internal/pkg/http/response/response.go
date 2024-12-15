package response

import (
	"log/slog"
	"net/http"
)

func RenderError(w http.ResponseWriter, err *HTTPError) {
	slog.Error(err.Error(), slog.String("error", err.Err.Error()))
	http.Error(w, err.Error(), err.Code)
}
