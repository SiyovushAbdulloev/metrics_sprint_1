package metric

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	"testing"
)

func BenchmarkUseCase_StoreMetric(b *testing.B) {
	db := memory.NewDB([]entity.Metrics{})
	repo := memory.NewMetricRepository(db)
	uc := New(repo)

	metric := entity.Metrics{
		ID:    "bench_metric",
		MType: entity.Gauge,
		Value: ptrFloat64(123.456),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.StoreMetric(metric)
	}
}

func BenchmarkUseCase_StoreAll(b *testing.B) {
	db := memory.NewDB([]entity.Metrics{})
	repo := memory.NewMetricRepository(db)
	uc := New(repo)

	metrics := []entity.Metrics{
		{
			ID:    "bench_metric1",
			MType: entity.Gauge,
			Value: ptrFloat64(1.23),
		},
		{
			ID:    "bench_metric2",
			MType: entity.Counter,
			Delta: ptrInt64(42),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uc.StoreAll(metrics)
	}
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}
