package repository

import "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"

type MetricRepository interface {
	StoreMetric(metric entity.Metrics) (entity.Metrics, error)
	StoreAll(metrics []entity.Metrics) error
	GetMetric(metric entity.Metrics) (entity.Metrics, error)
	GetMetrics() ([]entity.Metrics, error)
	Check() error
}
