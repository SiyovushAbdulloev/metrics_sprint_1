package repository

import "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"

type MetricRepository interface {
	StoreMetric(metric entity.Metrics) entity.Metrics
	StoreAll(metrics []entity.Metrics) bool
	GetMetric(metric entity.Metrics) (entity.Metrics, bool)
	GetMetrics() []entity.Metrics
}
