package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/alibek2024/FundGo/internal/delivery/helpers"
)

type ErrorMiddleware struct {
	logger *slog.Logger
}

func NewErrorMiddleware(l *slog.Logger) ErrorMiddleware {
	return ErrorMiddleware{
		logger: l,
	}
}

var (
	recoverError = errors.New("An unexpected error has occurred on the server")
)

func (e *ErrorMiddleware) ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				e.logger.Error("panic recovered",
					slog.String("panic", fmt.Sprint(err)),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("stack", string(debug.Stack())),
				)

				helpers.RespondWithError(w, helpers.InternalServerError, recoverError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
