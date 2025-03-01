package models

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Name  string
	Type  MetricType
	Value any
}

func NewMetric(name string, metricType MetricType, value any) *Metric {
	return &Metric{
		Name:  name,
		Type:  metricType,
		Value: value,
	}
}
