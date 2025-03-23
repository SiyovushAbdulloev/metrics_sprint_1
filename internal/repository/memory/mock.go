package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
)

type MockMemStorage struct {
	data []entity.Metric
}

type MockMetricRepository struct {
	DB MockMemStorage
}

func NewMockMetricRepository(db MockMemStorage) MockMetricRepository {
	return MockMetricRepository{
		DB: db,
	}
}

func NewMockDB(data []entity.Metric) MockMemStorage {
	return MockMemStorage{
		data: data,
	}
}

func (mms MockMetricRepository) StoreMetric(metric entity.Metric) bool {
	switch metric.Type {
	case entity.Gauge:
		data := mms.DB.data
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

		mms.DB.data = data
		return true
	case entity.Counter:
		data := mms.DB.data

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

		mms.DB.data = data
		return true
	default:
		return false
	}
}

func (mms MockMetricRepository) GetMetric(metricType string, metricName string) (entity.Metric, bool) {
	for _, m := range mms.DB.data {
		if m.Name == metricName && string(m.Type) == metricType {
			return m, true
		}
	}

	return entity.Metric{}, false
}

func (mms MockMetricRepository) GetMetrics() []entity.Metric {
	return mms.DB.data
}
