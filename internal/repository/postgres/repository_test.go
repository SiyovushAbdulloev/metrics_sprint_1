package postgres_test

import (
	"context"
	"flag"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/postgres"
	pgpkg "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	testRepo postgres.MetricRepository
	db       *pgpkg.Postgres
	dsn      string
)

func TestMain(m *testing.M) {
	flag.StringVar(&dsn, "dsn", "postgres://postgres:password@localhost:5432/metrics", "PostgreSQL DSN for testing")
	flag.Parse()

	if dsn == "" {
		panic("missing required flag: -dsn")
	}

	var err error
	db, err = pgpkg.New(dsn)
	if err != nil {
		panic(err)
	}

	testRepo = postgres.NewMetricRepository(db)
	code := m.Run()

	db.Close()
	os.Exit(code)
}

func cleanupMetricsTable(t *testing.T) {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM metrics")
	assert.NoError(t, err, "failed to clean metrics table")
}

func TestStoreAndGetMetric(t *testing.T) {
	metric := entity.Metrics{
		ID:    "test_metric",
		MType: "gauge",
		Value: ptrFloat64(42.0),
	}

	// Сохраняем метрику
	stored, err := testRepo.StoreMetric(metric)
	assert.NoError(t, err, "StoreMetric вернул ошибку")

	// Проверяем, что вернулась та же метрика
	assert.Equal(t, metric.ID, stored.ID)
	assert.Equal(t, metric.MType, stored.MType)

	// Получаем метрику
	got, err := testRepo.GetMetric(metric)
	assert.NoError(t, err, "GetMetric вернул ошибку")

	// Проверки содержимого
	assert.Equal(t, metric.ID, got.ID)
	assert.Equal(t, metric.MType, got.MType)
	require.NotNil(t, got.Value, "Значение Value должно быть не nil")
	assert.Equal(t, *metric.Value, *got.Value)
}

func TestStoreAllAndGetMetrics(t *testing.T) {
	cleanupMetricsTable(t)

	metrics := []entity.Metrics{
		{ID: "bulk1", MType: "counter", Delta: ptrInt64(3)},
		{ID: "bulk2", MType: "gauge", Value: ptrFloat64(99.9)},
	}

	err := testRepo.StoreAll(metrics)
	assert.NoError(t, err)

	all, err := testRepo.GetMetrics()
	assert.NoError(t, err)

	var found int
	for _, m := range all {
		if m.ID == "bulk1" || m.ID == "bulk2" {
			found++
		}
	}
	assert.Equal(t, 2, found)
}

func TestUpdateAll(t *testing.T) {
	cleanupMetricsTable(t)

	metrics := []entity.Metrics{
		{ID: "update1", MType: "counter", Delta: ptrInt64(1)},
		{ID: "update2", MType: "gauge", Value: ptrFloat64(10.5)},
	}
	_ = testRepo.StoreAll(metrics)

	updates := []entity.Metrics{
		{ID: "update1", MType: "counter", Delta: ptrInt64(4)},
		{ID: "update2", MType: "gauge", Value: ptrFloat64(12.3)},
	}
	err := testRepo.UpdateAll(updates)
	assert.NoError(t, err)

	got1, _ := testRepo.GetMetric(entity.Metrics{ID: "update1", MType: "counter"})
	got2, _ := testRepo.GetMetric(entity.Metrics{ID: "update2", MType: "gauge"})

	assert.Equal(t, int64(5), *got1.Delta) // 1 + 4
	assert.Equal(t, 12.3, *got2.Value)
}

func TestCheck(t *testing.T) {
	err := testRepo.Check()
	assert.NoError(t, err)
}

func ptrFloat64(v float64) *float64 { return &v }
func ptrInt64(v int64) *int64       { return &v }
