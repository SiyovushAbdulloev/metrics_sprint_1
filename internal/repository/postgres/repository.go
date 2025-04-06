package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
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

func (repo MetricRepository) StoreMetric(metric entity.Metrics) (entity.Metrics, error) {
	query := repo.DB.Builder.Insert("metrics").
		Columns("id", "type", "delta", "value").
		Values(metric.ID, metric.MType, metric.Delta, metric.Value).
		Suffix("ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, delta = EXCLUDED.delta, value = EXCLUDED.value")

	sql, args, err := query.ToSql()
	if err != nil {
		return entity.Metrics{}, err
	}

	_, err = repo.DB.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return entity.Metrics{}, err
	}

	return metric, nil
}

func (repo MetricRepository) StoreAll(metrics []entity.Metrics) error {
	query := repo.DB.Builder.Insert("metrics").
		Columns("id", "type", "delta", "value")

	for _, metric := range metrics {
		query = query.Values(metric.ID, metric.MType, metric.Delta, metric.Value)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = repo.DB.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo MetricRepository) GetMetric(metric entity.Metrics) (entity.Metrics, error) {
	query := repo.DB.Builder.
		Select("id, type, delta, value").
		From("metrics").
		Where(squirrel.Eq{"id": metric.ID, "type": metric.MType})

	sql, args, err := query.ToSql()
	if err != nil {
		return entity.Metrics{}, err
	}

	var m entity.Metrics
	err = repo.DB.Pool.QueryRow(context.Background(), sql, args...).Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return entity.Metrics{}, err
	}

	return m, nil
}

func (repo MetricRepository) GetMetrics() ([]entity.Metrics, error) {
	query := repo.DB.Builder.
		Select("id, type, delta, value").
		From("metrics")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.DB.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var metrics []entity.Metrics
	for rows.Next() {
		var m entity.Metrics
		if err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (repo MetricRepository) Check() error {
	return repo.DB.Pool.Ping(context.Background())
}
