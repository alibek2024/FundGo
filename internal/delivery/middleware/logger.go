package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *slog.Logger
}

func NewLoggerMiddleware(l *slog.Logger) LoggerMiddleware {
	return LoggerMiddleware{
		logger: l,
	}
}

type responseWriterDelegator struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriterDelegator) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterDelegator) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

func (l *LoggerMiddleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriterDelegator{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		level := slog.LevelInfo
		if wrapped.statusCode >= 500 {
			level = slog.LevelError
		} else if wrapped.statusCode >= 400 {
			level = slog.LevelWarn
		}

		l.logger.LogAttrs(
			r.Context(),
			level,
			"HTTP request handled",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", duration),
			slog.Int64("bytes", wrapped.written),
			slog.String("ip", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}
