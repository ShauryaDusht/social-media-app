package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counter for total requests
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Histogram for request duration
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Gauge for active requests
	httpActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
	)

	// Counter for requests per minute (calculated in Grafana)
	httpRequestsPerMinute = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_per_minute_total",
			Help: "Total HTTP requests for rate calculation",
		},
		[]string{"endpoint"},
	)
)

// collects HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		httpActiveRequests.Inc() // Increment active requests

		c.Next()

		httpActiveRequests.Dec() // Decrement active requests

		duration := time.Since(start).Seconds()

		endpoint := getCleanPath(c.FullPath())
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		httpRequestsPerMinute.WithLabelValues(endpoint).Inc()
	}
}

// getCleanPath removes path parameters for cleaner metrics
func getCleanPath(path string) string {
	if path == "" {
		return "unknown"
	}
	return path
}
