package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// ResponseCounter ...
	ResponseCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "endpoint"},
	)

	// ErrorCounter ...
	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_total",
			Help: "Total Error counts",
		},
		[]string{"method", "endpoint"},
	)

	// ResponseLatency ...
	ResponseLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "response_latency_millisecond",
			Help: "Response latency (millisecond)",
		},
		[]string{"method", "endpoint"},
	)

	// MetricHandler ...
	MetricHandler = gin.WrapH(promhttp.Handler())
)

func init() {
	fmt.Println("initing prometheus....")
	prometheus.MustRegister(ResponseCounter)
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(ResponseLatency)
}

// Prometheus ...
func Prometheus(pathBlacklist ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(pathBlacklist); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range pathBlacklist {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		endPoint := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key != "" {
				endPoint = strings.Replace(endPoint, p.Value, ":"+p.Key, 1)
			}
		}
		method := c.Request.Method
		start := time.Now()

		c.Next()

		if _, ok := skip[endPoint]; !ok {
			if length := len(c.Errors); length > 0 {
				ErrorCounter.WithLabelValues(method, endPoint).Add(float64(length))
			} else {
				statusCode := strconv.Itoa(c.Writer.Status())
				elapsed := float64(time.Since(start)) / float64(time.Millisecond)

				ResponseCounter.WithLabelValues(statusCode, method, endPoint).Inc()
				ResponseLatency.WithLabelValues(method, endPoint).Observe(elapsed)
			}
		}
	}
}
