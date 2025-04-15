package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	err "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
	"strings"
)

type MockMemStorage struct {
	data []entity.Metrics
}

type MockMetricRepository struct {
	DB MockMemStorage
}

func NewMockMetricRepository(db MockMemStorage) MockMetricRepository {
	return MockMetricRepository{
		DB: db,
	}
}

func NewMockDB(data []entity.Metrics) MockMemStorage {
	return MockMemStorage{
		data: data,
	}
}

func (mms MockMetricRepository) StoreMetric(metric entity.Metrics) (entity.Metrics, error) {
	data := mms.DB.data

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

		mms.DB.data = data
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

		mms.DB.data = data
		return metric, nil
	default:
		return entity.Metrics{}, err.ErrNotFound
	}
}

func (mms MockMetricRepository) GetMetric(metric entity.Metrics) (entity.Metrics, error) {
	for _, m := range mms.DB.data {
		if m.ID == strings.ToLower(metric.ID) && m.MType == metric.MType {
			return m, nil
		}
	}

	return entity.Metrics{}, nil
}

func (mms MockMetricRepository) GetMetrics() ([]entity.Metrics, error) {
	return mms.DB.data, nil
}

func (mms MockMetricRepository) StoreAll(metrics []entity.Metrics) error {
	mms.DB.data = metrics

	return nil
}

func (mms MockMetricRepository) Check() error {
	return nil
}

func (mms MockMetricRepository) UpdateAll(metrics []entity.Metrics) error {
	return nil
}
