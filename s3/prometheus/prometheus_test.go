package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/stretchr/testify/assert"
)

func TestCollectors(t *testing.T) {
	colls := []prometheus.Collector{
		collectors.NewBuildInfoCollector(),
		requestCounter,
		latencyHistogram,
	}

	for _, coll := range colls {
		err := prometheus.Register(coll)
		are, ok := err.(prometheus.AlreadyRegisteredError)
		if ok {
			assert.Equal(t, coll, are.ExistingCollector)
		}
		assert.True(t, ok)
	}
}
