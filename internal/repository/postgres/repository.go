// Package postgres содержит реализацию репозитория метрик на основе PostgreSQL.
package postgres

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"log"
)

// MetricRepository предоставляет доступ к хранилищу метрик в PostgreSQL.
type MetricRepository struct {
	DB *postgres.Postgres
}

// NewMetricRepository создаёт новый экземпляр MetricRepository.
func NewMetricRepository(db *postgres.Postgres) MetricRepository {
	return MetricRepository{
		DB: db,
	}
}

// StoreMetric сохраняет или обновляет одну метрику в базе данных.
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

// StoreAll сохраняет список метрик в базу данных.
// Используется для пакетной вставки.
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

// GetMetric возвращает одну метрику по ID и типу из базы данных.
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

// GetMetrics возвращает все метрики из базы данных.
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
		if err2 := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err2 != nil {
			return nil, err2
		}
		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

// Check выполняет проверку соединения с базой данных.
func (repo MetricRepository) Check() error {
	return repo.DB.Pool.Ping(context.Background())
}

// UpdateAll обновляет значения существующих метрик или вставляет новые,
// используя транзакцию и UPSERT-логику.
func (repo MetricRepository) UpdateAll(metrics []entity.Metrics) error {
	ctx := context.Background()

	tx, err := repo.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		log.Printf("Error defer: %v\n", err)
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	query := repo.DB.Builder.Insert("metrics").
		Columns("id", "type", "delta", "value")

	for _, metric := range metrics {
		var delta interface{}
		var value interface{}

		if metric.Delta != nil {
			oldMetric, err2 := repo.GetMetric(metric)

			if err2 != nil && !errors.Is(err2, pgx.ErrNoRows) {
				return err2
			}
			if errors.Is(err2, pgx.ErrNoRows) {
				delta = *metric.Delta
			} else {
				delta = *oldMetric.Delta + *metric.Delta
			}
		} else {
			delta = nil
		}

		if metric.Value != nil {
			value = *metric.Value
		} else {
			value = nil
		}

		query = query.Values(metric.ID, metric.MType, delta, value)
	}

	query = query.Suffix("ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)

	return err
}
