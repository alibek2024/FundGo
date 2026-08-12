package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/alibek2024/FundGo/internal/delivery"
)

type errorMiddleware struct {
	logger *slog.Logger
}

func NewErrorMiddleware(l *slog.Logger) errorMiddleware {
	return errorMiddleware{
		logger: l,
	}
}

var (
	recoverError = errors.New("An unexpected error has occurred on the server")
)

func (e *errorMiddleware) ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				e.logger.Error("panic recovered",
					slog.String("panic", fmt.Sprint(err)),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("stack", string(debug.Stack())),
				)

				RespondWithError(w, delivery.InternalServerError, recoverError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	errMsg := "unknown error"
    if err != nil {
        errMsg = err.Error()
    }
	Respond(w, code, map[string]string{"error": errMsg})
}

func Respond(w http.ResponseWriter, code int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if message != nil {
		json.NewEncoder(w).Encode(message)
	}
}