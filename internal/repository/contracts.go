package repository

import "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"

type MetricRepository interface {
	StoreMetric(metric entity.Metric) bool
	GetMetric(metricType string, metricName string) (entity.Metric, bool)
	GetMetrics() []entity.Metric
}
