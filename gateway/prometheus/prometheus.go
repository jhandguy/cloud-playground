package prometheus

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/zap"
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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func CollectMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/monitoring/") {
			next.ServeHTTP(w, r)
			return
		}

		route := mux.CurrentRoute(r)
		path, err := route.GetPathTemplate()
		if err != nil {
			zap.S().Errorw("failed to get path template", "error", err)
		}

		rw := &responseWriter{
			ResponseWriter: w,
		}
		startTime := time.Now()
		next.ServeHTTP(rw, r)

		success := false
		if rw.statusCode < http.StatusBadRequest {
			success = true
		}

		requestCounter.
			WithLabelValues(path, r.Method, strconv.FormatBool(success)).
			Inc()

		latencyHistogram.
			WithLabelValues(path, r.Method).
			Observe(time.Since(startTime).Seconds())
	})
}
