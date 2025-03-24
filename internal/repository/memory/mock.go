package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
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

func (mms MockMetricRepository) StoreMetric(metric entity.Metrics) entity.Metrics {
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

		mms.DB.data = data
		return metric
	default:
		return entity.Metrics{}
	}
}

func (mms MockMetricRepository) GetMetric(metric entity.Metrics) (entity.Metrics, bool) {
	for _, m := range mms.DB.data {
		if m.ID == strings.ToLower(metric.ID) && m.MType == metric.MType {
			return m, true
		}
	}

	return entity.Metrics{}, false
}

func (mms MockMetricRepository) GetMetrics() []entity.Metrics {
	return mms.DB.data
}

func (mms MockMetricRepository) StoreAll(metrics []entity.Metrics) bool {
	mms.DB.data = metrics

	return true
}
