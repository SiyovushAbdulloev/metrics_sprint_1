package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
)

type MemStorage struct {
	data []models.Metric
}

type MetricStorage struct {
	DB MemStorage
}

func NewMetricStorage(db MemStorage) MetricStorage {
	return MetricStorage{
		DB: db,
	}
}

func NewDB(data []models.Metric) MemStorage {
	return MemStorage{
		data: data,
	}
}

func (ms *MetricStorage) StoreMetric(metric models.Metric) bool {
	switch metric.Type {
	case models.Gauge:
		data := ms.DB.data
		if len(data) == 0 {
			data = append(data, metric)
		} else {
			notFound := true
			for k, v := range data {
				if v.Name == metric.Name {
					notFound = false
					data[k] = metric
				}
			}

			if notFound {
				data = append(data, metric)
			}
		}

		ms.DB.data = data
		return true
	case models.Counter:
		data := ms.DB.data

		if len(data) == 0 {
			data = append(data, metric)
		} else {
			notFound := true
			for k, v := range data {
				if v.Name == metric.Name {
					notFound = false
					newValue := metric.Value.(int64)
					previousValue := v.Value.(int64)
					data[k] = models.Metric{
						Name:  metric.Name,
						Type:  metric.Type,
						Value: newValue + previousValue,
					}
				}
			}

			if notFound {
				data = append(data, metric)
			}
		}

		ms.DB.data = data
		return true
	default:
		return false
	}
}
