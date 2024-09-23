package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	PaymentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_total",
			Help: "Total number of payments processed",
		},
		[]string{"status"},
	)

	PaymentAmount = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_amount",
			Help:    "Distribution of payment amounts",
			Buckets: prometheus.ExponentialBuckets(10, 2, 10), // 10, 20, 40, ..., 5120
		},
		[]string{"currency"},
	)

	RefundTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "refund_total",
			Help: "Total number of refunds processed",
		},
		[]string{"status"},
	)

	RefundAmount = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "refund_amount",
			Help:    "Distribution of refund amounts",
			Buckets: prometheus.ExponentialBuckets(10, 2, 10),
		},
		[]string{"currency"},
	)

	APIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	DatabaseQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10), // 1ms, 2ms, 4ms, ..., 512ms
		},
		[]string{"query_type"},
	)

	AuthenticationAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authentication_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"status"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(PaymentTotal)
	prometheus.MustRegister(PaymentAmount)
	prometheus.MustRegister(RefundTotal)
	prometheus.MustRegister(RefundAmount)
	prometheus.MustRegister(APIRequestDuration)
	prometheus.MustRegister(DatabaseQueryDuration)
	prometheus.MustRegister(AuthenticationAttempts)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
