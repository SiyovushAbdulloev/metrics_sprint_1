package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
)

type MockMemStorage struct {
	data []models.Metric
}

type MockMetricStorage struct {
	DB MockMemStorage
}

func NewMockMetricStorage(db MockMemStorage) MockMetricStorage {
	return MockMetricStorage{
		DB: db,
	}
}

func NewMockDB(data []models.Metric) MockMemStorage {
	return MockMemStorage{
		data: data,
	}
}

func (mms *MockMetricStorage) StoreMetric(metric models.Metric) bool {
	switch metric.Type {
	case models.Gauge:
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
	case models.Counter:
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

		mms.DB.data = data
		return true
	default:
		return false
	}
}

func (mms MockMetricStorage) GetMetric(metricType string, metricName string) (models.Metric, bool) {
	for _, m := range mms.DB.data {
		if m.Name == metricName && string(m.Type) == metricType {
			return m, true
		}
	}

	return models.Metric{}, false
}

func (mms MockMetricStorage) GetMetrics() []models.Metric {
	return mms.DB.data
}
