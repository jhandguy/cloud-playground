package prometheus

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_gateway_requests_count",
			Help: "Request counter per path and method",
		},
		[]string{"path", "method", "success"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "devops_playground_gateway_requests_latency",
			Help: "Request latency histogram per path and method",
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(latencyHistogram)
}

func CollectMetrics(c *gin.Context) {
	if strings.Contains(c.Request.RequestURI, "/monitoring/") {
		c.Next()
		return
	}

	startTime := time.Now()
	c.Next()

	success := false
	if c.Writer.Status() < http.StatusBadRequest {
		success = true
	}

	requestCounter.
		WithLabelValues(c.FullPath(), c.Request.Method, strconv.FormatBool(success)).
		Inc()

	latencyHistogram.
		WithLabelValues(c.FullPath(), c.Request.Method).
		Observe(time.Since(startTime).Seconds())
}
