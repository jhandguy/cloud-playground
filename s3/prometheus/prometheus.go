package prometheus

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"google.golang.org/grpc"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloud_playground_s3_requests_count",
			Help: "Request counter per method",
		},
		[]string{"method", "success"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cloud_playground_s3_requests_latency",
			Help: "Request latency histogram per method",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(latencyHistogram)
}

func CollectMetrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	startTime := time.Now()
	res, err := handler(ctx, req)

	success := false
	if err == nil {
		success = true
	}

	requestCounter.
		WithLabelValues(info.FullMethod, strconv.FormatBool(success)).
		Inc()

	latencyHistogram.
		WithLabelValues(info.FullMethod).
		Observe(time.Since(startTime).Seconds())

	return res, err
}
