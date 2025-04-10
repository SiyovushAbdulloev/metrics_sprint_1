package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	err "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
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

func (ms MetricRepository) StoreMetric(metric entity.Metrics) (entity.Metrics, error) {
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
		return metric, nil
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
		return metric, nil
	default:
		return entity.Metrics{}, err.ErrInvalidType
	}
}

func (ms MetricRepository) GetMetric(metric entity.Metrics) (entity.Metrics, error) {
	for _, m := range ms.DB.data {
		if m.ID == metric.ID && m.MType == metric.MType {
			return m, nil
		}
	}

	return entity.Metrics{}, err.ErrNotFound
}

func (ms MetricRepository) GetMetrics() ([]entity.Metrics, error) {
	return ms.DB.data, nil
}

func (ms MetricRepository) StoreAll(metrics []entity.Metrics) error {
	ms.DB.data = metrics

	return nil
}

func (ms MetricRepository) Check() error {
	return nil
}

func (ms MetricRepository) UpdateAll(metrics []entity.Metrics) error {
	return nil
}
