package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func ServerError(w http.ResponseWriter, msg string, err error) {
	const errorText = "An error occurred."

	slog.Error(msg, "Error:", err)
	http.Error(w, errorText, http.StatusInternalServerError)
}
