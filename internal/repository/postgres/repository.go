package postgres

import (
	"context"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
)

type MetricRepository struct {
	DB *postgres.Postgres
}

func NewMetricRepository(db *postgres.Postgres) MetricRepository {
	return MetricRepository{
		DB: db,
	}
}

func (repo MetricRepository) StoreMetric(metric entity.Metrics) entity.Metrics {
	return entity.Metrics{}
}

func (repo MetricRepository) StoreAll(metrics []entity.Metrics) bool {
	return true
}

func (repo MetricRepository) GetMetric(metric entity.Metrics) (entity.Metrics, bool) {
	return entity.Metrics{}, false
}

func (repo MetricRepository) GetMetrics() []entity.Metrics {
	return []entity.Metrics{}
}

func (repo MetricRepository) Check() error {
	return repo.DB.Pool.Ping(context.Background())
}
