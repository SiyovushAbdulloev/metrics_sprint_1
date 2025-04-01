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

func (uc UseCase) StoreMetric(metric entity.Metric) bool {
	return uc.repo.StoreMetric(metric)
}

func (uc UseCase) GetMetric(metricType string, metricName string) (entity.Metric, bool) {
	return uc.repo.GetMetric(metricType, metricName)
}

func (uc UseCase) GetMetrics() []entity.Metric {
	return uc.repo.GetMetrics()
}
