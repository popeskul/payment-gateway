package middleware

import (
	"net/http"
	"time"

	"github.com/popeskul/payment-gateway/internal/infrastructure/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := NewStatusResponseWriter(w)
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		metrics.APIRequestDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(ww.status)).Observe(duration)
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{w, http.StatusOK}
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
