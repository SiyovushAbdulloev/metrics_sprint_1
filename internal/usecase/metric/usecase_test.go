package metric

import (
	metric "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/postgres_metric"
	"testing"

	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	"github.com/stretchr/testify/assert"
)

func setupUseCase() *metric.UseCase {
	db := memory.NewDB([]entity.Metrics{})
	repo := memory.NewMetricRepository(db)
	return metric.New(repo)
}

func TestUseCase_StoreAndGetMetric(t *testing.T) {
	uc := setupUseCase()

	value := 123.45
	metricData := entity.Metrics{
		ID:    "metric1",
		MType: entity.Gauge,
		Value: &value,
	}

	stored, err := uc.StoreMetric(metricData)
	assert.NoError(t, err)
	assert.Equal(t, "metric1", stored.ID)

	got, err := uc.GetMetric(metricData)
	assert.NoError(t, err)
	assert.Equal(t, metricData.ID, got.ID)
	assert.Equal(t, metricData.MType, got.MType)
	assert.NotNil(t, got.Value)
	assert.Equal(t, *metricData.Value, *got.Value)
}

func TestUseCase_StoreAllAndGetMetrics(t *testing.T) {
	uc := setupUseCase()

	metrics := []entity.Metrics{
		{ID: "batch1", MType: entity.Counter, Delta: ptrInt64(5)},
		{ID: "batch2", MType: entity.Gauge, Value: ptrFloat64(42.0)},
	}

	err := uc.StoreAll(metrics)
	assert.NoError(t, err)

	all, err := uc.GetMetrics()
	assert.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestUseCase_Check(t *testing.T) {
	uc := setupUseCase()

	err := uc.Check()
	assert.NoError(t, err)
}

func TestUseCase_UpdateAll(t *testing.T) {
	uc := setupUseCase()

	err := uc.UpdateAll([]entity.Metrics{
		{ID: "dummy", MType: entity.Gauge, Value: ptrFloat64(1)},
	})
	assert.NoError(t, err)
}

func ptrFloat64(v float64) *float64 { return &v }
func ptrInt64(v int64) *int64       { return &v }
