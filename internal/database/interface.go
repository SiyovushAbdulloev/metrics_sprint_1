package database

import "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"

type MetricStorage interface {
	StoreMetric(metric models.Metric) bool
}
