package opentelemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTracer(t *testing.T) {
	tracer := GetTracer("test")
	assert.Nil(t, tracer)
}

func TestGetTraceID(t *testing.T) {
	traceID := GetTraceID(context.Background())
	assert.Empty(t, traceID)
}

func TestStopTracing(t *testing.T) {
	err := StopTracing(context.Background())
	assert.Nil(t, err)
}
