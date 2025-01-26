package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type Middleware func(http.Handler) http.Handler

func RecoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(
					"Panic occurred.",
					"reason", err,
					"request", fmt.Sprint(r),
					"stack_trace", string(debug.Stack()),
				)
				status := http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
