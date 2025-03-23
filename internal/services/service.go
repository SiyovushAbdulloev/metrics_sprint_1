package services

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/database"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
)

type MetricService struct {
	Storage database.MetricStorage
}

func NewMetricService(storage database.MetricStorage) MetricService {
	return MetricService{
		Storage: storage,
	}
}

func (ms MetricService) StoreMetric(metric models.Metric) bool {
	return ms.Storage.StoreMetric(metric)
}

func (ms MetricService) GetMetric(metricType string, metricName string) (models.Metric, bool) {
	return ms.Storage.GetMetric(metricType, metricName)
}

func (ms MetricService) GetMetrics() []models.Metric {
	return ms.Storage.GetMetrics()
}
