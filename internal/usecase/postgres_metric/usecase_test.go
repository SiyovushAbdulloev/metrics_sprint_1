//go:build integration
// +build integration

package metric_test

import (
	"context"
	"flag"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/postgres"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/postgres_metric"
	postgres2 "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	metricUC *metric.UseCase
	db       *postgres2.Postgres
)

func TestMain(m *testing.M) {
	dsn := flag.String("dsn", "postgres://postgres:password@localhost:5432/metrics", "PostgreSQL DSN")
	flag.Parse()

	if *dsn == "" {
		panic("DSN must be provided with -dsn flag")
	}

	var err error
	db, err = postgres2.New(*dsn)
	if err != nil {
		panic(err)
	}

	repo := postgres.NewMetricRepository(db)
	metricUC = metric.New(repo)

	code := m.Run()

	db.Close()
}

func cleanupMetricsTable(t *testing.T) {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM metrics")
	assert.NoError(t, err, "failed to clean metrics table")
}

func TestStoreAndGetMetric(t *testing.T) {
	cleanupMetricsTable(t)

	delta := int64(123)
	metric := entity.Metrics{
		ID:    "test_metric_uc",
		MType: entity.Counter,
		Delta: &delta,
	}

	stored, err := metricUC.StoreMetric(metric)
	assert.NoError(t, err)
	assert.Equal(t, metric.ID, stored.ID)

	retrieved, err := metricUC.GetMetric(metric)
	assert.NoError(t, err)
	assert.Equal(t, metric.ID, retrieved.ID)
	assert.Equal(t, *metric.Delta, *retrieved.Delta)
}

func TestStoreAllAndGetMetrics(t *testing.T) {
	cleanupMetricsTable(t)

	value := 42.42
	metrics := []entity.Metrics{
		{
			ID:    "batch_metric_1",
			MType: entity.Gauge,
			Value: &value,
		},
		{
			ID:    "batch_metric_2",
			MType: entity.Gauge,
			Value: &value,
		},
	}

	err := metricUC.StoreAll(metrics)
	assert.NoError(t, err)

	allMetrics, err := metricUC.GetMetrics()
	assert.NoError(t, err)

	found := 0
	for _, m := range allMetrics {
		if m.ID == "batch_metric_1" || m.ID == "batch_metric_2" {
			found++
		}
	}
	assert.Equal(t, 2, found)
}

func TestCheckConnection(t *testing.T) {
	cleanupMetricsTable(t)

	err := metricUC.Check()
	assert.NoError(t, err)
}
