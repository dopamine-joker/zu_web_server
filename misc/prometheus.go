package misc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	NAMESPACE = "gin"
)

var (
	Handler     = "handler"
	Method      = "method"
	Code        = "code"
	ServiceName = "service_name"
)

type PrometheusMonitor struct {
	ServiceName        string //服务名
	APIRequestsCounter *prometheus.CounterVec
	RequestDuration    *prometheus.HistogramVec
	RequestSize        *prometheus.HistogramVec
	ResponseSize       *prometheus.HistogramVec
}

func NewPrometheusMonitor(nameSpace, serviceName string) *PrometheusMonitor {
	labels := []string{Handler, Method, Code, ServiceName}
	APIRequestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: nameSpace,
			Name:      "http_request_count",
			Help:      "A counter for requests to the wrapped handler",
		},
		labels,
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Name:      "http_request_duration_seconds",
			Help:      "A histogram of latencies for requests.",
		},
		labels,
	)
	requestSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Name:      "http_request_size_bytes",
			Help:      "A histogram of request sizes for requests.",
		},
		labels,
	)
	responseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: nameSpace,
			Name:      "http_response_size_bytes",
			Help:      "A histogram of response sizes for requests.",
		},
		labels,
	)
	prometheus.MustRegister(APIRequestsCounter, requestDuration, requestSize, responseSize)
	return &PrometheusMonitor{
		ServiceName:        serviceName,
		APIRequestsCounter: APIRequestsCounter,
		RequestDuration:    requestDuration,
		RequestSize:        requestSize,
		ResponseSize:       responseSize,
	}
}

func (m *PrometheusMonitor) PromMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		reqSize := computeApproximateRequestSize(c.Request)
		c.Next()
		duration := time.Since(start)
		statusCode := fmt.Sprintf("%d", c.Writer.Status())
		promLabels := prometheus.Labels{Handler: path, Method: c.Request.Method, Code: statusCode, ServiceName: m.ServiceName}
		m.APIRequestsCounter.With(promLabels).Inc()
		m.RequestDuration.With(promLabels).Observe(duration.Seconds())
		m.RequestSize.With(promLabels).Observe(float64(reqSize))
		m.ResponseSize.With(promLabels).Observe(float64(c.Writer.Size()))
	}
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

func StartMonitor(r *gin.Engine) {
	r.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})
}
