package memory

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	errPkg "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ptrFloat64(v float64) *float64 { return &v }
func ptrInt64(v int64) *int64       { return &v }

func TestMemoryMetricRepository_StoreAndGetGauge(t *testing.T) {
	repo := NewMetricRepository(NewDB(nil))
	metric := entity.Metrics{ID: "gauge1", MType: entity.Gauge, Value: ptrFloat64(99.9)}

	stored, err := repo.StoreMetric(metric)
	assert.NoError(t, err)
	assert.Equal(t, metric, stored)

	fetched, err := repo.GetMetric(metric)
	assert.NoError(t, err)
	assert.Equal(t, metric, fetched)
}

func TestMemoryMetricRepository_StoreAndGetCounter(t *testing.T) {
	repo := NewMetricRepository(NewDB(nil))
	metric := entity.Metrics{ID: "count1", MType: entity.Counter, Delta: ptrInt64(10)}

	_, err := repo.StoreMetric(metric)
	assert.NoError(t, err)

	// Add more delta to same metric
	metric2 := entity.Metrics{ID: "count1", MType: entity.Counter, Delta: ptrInt64(5)}
	_, err = repo.StoreMetric(metric2)
	assert.NoError(t, err)

	expected := int64(15)
	fetched, err := repo.GetMetric(entity.Metrics{ID: "count1", MType: entity.Counter})
	assert.NoError(t, err)
	assert.Equal(t, expected, *fetched.Delta)
}

func TestMemoryMetricRepository_GetNotFound(t *testing.T) {
	repo := NewMetricRepository(NewDB(nil))
	_, err := repo.GetMetric(entity.Metrics{ID: "missing", MType: entity.Gauge})
	assert.ErrorIs(t, err, errPkg.ErrNotFound)
}

func TestMemoryMetricRepository_StoreAllAndGetMetrics(t *testing.T) {
	repo := NewMetricRepository(NewDB(nil))
	metrics := []entity.Metrics{
		{ID: "m1", MType: entity.Gauge, Value: ptrFloat64(1.1)},
		{ID: "m2", MType: entity.Counter, Delta: ptrInt64(2)},
	}
	err := repo.StoreAll(metrics)
	assert.NoError(t, err)

	all, err := repo.GetMetrics()
	assert.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestMemoryMetricRepository_InvalidType(t *testing.T) {
	repo := NewMetricRepository(NewDB(nil))
	metric := entity.Metrics{ID: "bad", MType: "unknown"}
	_, err := repo.StoreMetric(metric)
	assert.ErrorIs(t, err, errPkg.ErrInvalidType)
}

func TestMemoryMetricRepository_CheckAndUpdateAll(t *testing.T) {
	repo := NewMetricRepository(NewDB(nil))
	err := repo.Check()
	assert.NoError(t, err)

	err = repo.UpdateAll(nil)
	assert.NoError(t, err) // current impl is no-op
}
