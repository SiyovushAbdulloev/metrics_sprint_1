package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	m := Metrics{}

	collectMetrics(&m)

	assert.NotEmpty(t, m.data)

	var pollCount int64
	for _, metric := range m.data {
		if metric.Name == "poll_count" {
			pollCount = metric.Value.(int64)
			break
		}
	}

	assert.Equal(t, int64(1), pollCount)

	collectMetrics(&m)
	for _, metric := range m.data {
		if metric.Name == "poll_count" {
			pollCount = metric.Value.(int64)
			break
		}
	}

	assert.Equal(t, int64(2), pollCount)
}
