package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestCollectors(t *testing.T) {
	collectors := []prometheus.Collector{
		prometheus.NewBuildInfoCollector(),
		totalReqCounter,
		successReqCounter,
		latencyHistogram,
	}

	for _, collector := range collectors {
		err := prometheus.Register(collector)
		are, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			assert.Equal(t, collector, are.ExistingCollector)
		}
		assert.True(t, ok)
	}
}
