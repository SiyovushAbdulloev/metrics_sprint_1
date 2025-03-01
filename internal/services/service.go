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
