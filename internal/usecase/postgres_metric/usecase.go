// Package metric реализует слой бизнес-логики для работы с метриками.
package metric

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository"
)

// UseCase представляет реализацию бизнес-логики для работы с метриками.
type UseCase struct {
	repo repository.MetricRepository
}

// New создаёт новый экземпляр UseCase с указанным репозиторием метрик.
func New(repo repository.MetricRepository) *UseCase {
	return &UseCase{repo: repo}
}

// StoreMetric сохраняет метрику (gauge или counter) в хранилище.
func (uc UseCase) StoreMetric(metric entity.Metrics) (entity.Metrics, error) {
	return uc.repo.StoreMetric(metric)
}

// GetMetric возвращает одну метрику по идентификатору и типу.
func (uc UseCase) GetMetric(metric entity.Metrics) (entity.Metrics, error) {
	return uc.repo.GetMetric(metric)
}

// GetMetrics возвращает все доступные метрики.
func (uc UseCase) GetMetrics() ([]entity.Metrics, error) {
	return uc.repo.GetMetrics()
}

// StoreAll сохраняет список метрик в хранилище.
func (uc UseCase) StoreAll(metrics []entity.Metrics) error {
	return uc.repo.StoreAll(metrics)
}

// Check проверяет доступность хранилища.
func (uc UseCase) Check() error {
	return uc.repo.Check()
}

// UpdateAll обновляет список метрик в хранилище.
func (uc UseCase) UpdateAll(metrics []entity.Metrics) error {
	return uc.repo.UpdateAll(metrics)
}
