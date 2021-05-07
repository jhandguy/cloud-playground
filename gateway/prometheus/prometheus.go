package prometheus

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

var (
	totalReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_gateway_requests_total",
			Help: "Total requests counter per path and method",
		},
		[]string{"path", "method", "deployment"},
	)

	successReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_gateway_requests_success",
			Help: "Successful requests counter per path and method",
		},
		[]string{"path", "method", "deployment"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "devops_playground_gateway_requests_latency",
			Help: "Requests latency histogram per path and method",
		},
		[]string{"path", "method", "deployment"},
	)
)

func init() {
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	prometheus.MustRegister(totalReqCounter)
	prometheus.MustRegister(successReqCounter)
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
		if strings.Contains(r.RequestURI, "/health") || strings.Contains(r.RequestURI, "/metrics") {
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

		totalReqCounter.
			WithLabelValues(path, r.Method, deployment).
			Inc()

		latencyHistogram.
			WithLabelValues(path, r.Method, deployment).
			Observe(time.Since(startTime).Seconds())

		if rw.statusCode < http.StatusBadRequest {
			successReqCounter.
				WithLabelValues(path, r.Method, deployment).
				Inc()
		}
	})
}
