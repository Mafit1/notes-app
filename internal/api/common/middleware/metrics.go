package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "total number of http requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

type MetricsMW struct{}

func NewMetricsMW() *MetricsMW {
	return &MetricsMW{}
}

func (m *MetricsMW) MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		method := c.Request().Method
		path := c.Path()

		defer func() {
			status := c.Response().Status
			duration := time.Since(start).Seconds()
			httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
			httpRequestDuration.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration)
		}()

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}
