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

func (sw *statusWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(log *zap.Logger, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Info("Incoming request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("addr", r.RemoteAddr),
		)

		wrappedWriter := &statusWriter{ResponseWriter: w}
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
