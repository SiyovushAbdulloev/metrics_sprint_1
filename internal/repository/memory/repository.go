package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
)

type MemStorage struct {
	data []entity.Metrics
}

type MetricRepository struct {
	DB *MemStorage
}

func NewMetricRepository(db MemStorage) MetricRepository {
	return MetricRepository{
		DB: &db,
	}
}

func NewDB(data []entity.Metrics) MemStorage {
	return MemStorage{
		data: data,
	}
}

func (ms MetricRepository) StoreMetric(metric entity.Metrics) entity.Metrics {
	data := ms.DB.data

	switch metric.MType {
	case entity.Gauge:
		if len(data) == 0 {
			data = append(data, metric)
		} else {
			notFound := true
			for k, v := range data {
				if v.ID == metric.ID {
					notFound = false
					data[k] = metric
				}
			}

			if notFound {
				data = append(data, metric)
			}
		}

		ms.DB.data = data
		return metric
	case entity.Counter:
		if len(data) == 0 {
			data = append(data, metric)
		} else {
			notFound := true
			for k, v := range data {
				if v.ID == metric.ID {
					notFound = false
					if metric.Delta != nil && v.Delta != nil {
						newValue := *metric.Delta + *v.Delta
						metric.Delta = &newValue
					}
					data[k] = metric
				}
			}

			if notFound {
				data = append(data, metric)
			}
		}

		ms.DB.data = data
		return metric
	default:
		return entity.Metrics{}
	}
}

func (ms MetricRepository) GetMetric(metric entity.Metrics) (entity.Metrics, bool) {
	for _, m := range ms.DB.data {
		if m.ID == metric.ID && m.MType == metric.MType {
			return m, true
		}
	}

	return entity.Metrics{}, false
}

func (ms MetricRepository) GetMetrics() []entity.Metrics {
	return ms.DB.data
}
