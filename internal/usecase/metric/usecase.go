package metric

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository"
)

type UseCase struct {
	repo repository.MetricRepository
}

func New(repo repository.MetricRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc UseCase) StoreMetric(metric entity.Metrics) entity.Metrics {
	return uc.repo.StoreMetric(metric)
}

func (uc UseCase) GetMetric(metric entity.Metrics) (entity.Metrics, bool) {
	return uc.repo.GetMetric(metric)
}

func (uc UseCase) GetMetrics() []entity.Metrics {
	return uc.repo.GetMetrics()
}

func (uc UseCase) StoreAll(metrics []entity.Metrics) bool {
	return uc.repo.StoreAll(metrics)
}
