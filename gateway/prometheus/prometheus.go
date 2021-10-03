package prometheus

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/viper"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_gateway_requests_count",
			Help: "Request counter per path and method",
		},
		[]string{"path", "method", "deployment", "success"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "devops_playground_gateway_requests_latency",
			Help: "Request latency histogram per path and method",
		},
		[]string{"path", "method", "deployment"},
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
	deployment := viper.GetString("gateway-deployment")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/monitoring/") {
			next.ServeHTTP(w, r)
			return
		}

		route := mux.CurrentRoute(r)
		path, err := route.GetPathTemplate()
		if err != nil {
			log.Printf("failed to get path template: %v", err)
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
			WithLabelValues(path, r.Method, deployment, strconv.FormatBool(success)).
			Inc()

		latencyHistogram.
			WithLabelValues(path, r.Method, deployment).
			Observe(time.Since(startTime).Seconds())
	})
}
