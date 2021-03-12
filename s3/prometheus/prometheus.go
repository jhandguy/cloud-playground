package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	totalReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_s3_requests_total",
			Help: "Total requests counter per method",
		},
		[]string{"method"},
	)

	successReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "devops_playground_s3_requests_success",
			Help: "Successful requests counter per method",
		},
		[]string{"method"},
	)

	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "devops_playground_s3_requests_latency",
			Help: "Requests latency histogram per method",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	prometheus.MustRegister(totalReqCounter)
	prometheus.MustRegister(successReqCounter)
	prometheus.MustRegister(latencyHistogram)
}

func CollectMetrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	startTime := time.Now()
	res, err := handler(ctx, req)

	totalReqCounter.
		WithLabelValues(info.FullMethod).
		Inc()

	latencyHistogram.
		WithLabelValues(info.FullMethod).
		Observe(time.Since(startTime).Seconds())

	if err == nil {
		successReqCounter.
			WithLabelValues(info.FullMethod).
			Inc()
	}

	return res, err
}
