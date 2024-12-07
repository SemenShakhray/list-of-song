package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.statusCode = code
	sw.ResponseWriter.WriteHeader(code)
}

func WithLogging(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("addr", r.RemoteAddr),
			)

			wrappedWriter := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)
			log.Info("Request completed",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", wrappedWriter.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}
