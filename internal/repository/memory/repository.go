package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
)

type MemStorage struct {
	data []entity.Metric
}

type MetricRepository struct {
	DB *MemStorage
}

func NewMetricRepository(db MemStorage) MetricRepository {
	return MetricRepository{
		DB: &db,
	}
}

func NewDB(data []entity.Metric) MemStorage {
	return MemStorage{
		data: data,
	}
}

func (ms MetricRepository) StoreMetric(metric entity.Metric) bool {
	switch metric.Type {
	case entity.Gauge:
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
	case entity.Counter:
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
					data[k] = entity.Metric{
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

func (ms MetricRepository) GetMetric(metricType string, metricName string) (entity.Metric, bool) {
	for _, m := range ms.DB.data {
		if m.Name == metricName && string(m.Type) == metricType {
			return m, true
		}
	}

	return entity.Metric{}, false
}

func (ms MetricRepository) GetMetrics() []entity.Metric {
	return ms.DB.data
}
