package metric

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase"
	error2 "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc usecase.MetricUseCase
	l  logger.Interface
}

func New(uc usecase.MetricUseCase, l logger.Interface) *Handler {
	return &Handler{
		uc: uc,
		l:  l,
	}
}

func (h *Handler) StoreMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")
	metricValue := ctx.Param("value")

	if metricType != string(entity.Gauge) && metricType != string(entity.Counter) {
		ctx.String(http.StatusBadRequest, error2.ErrInvalidType.Error())
		return
	}

	value, ok := h.validValue(metricType, metricValue)
	if !ok {
		ctx.String(http.StatusBadRequest, error2.ErrInvalidValue.Error())
		return
	}

	added := h.uc.StoreMetric(entity.Metric{
		Name:  metricName,
		Value: value,
		Type:  entity.MetricType(metricType),
	})

	if !added {
		ctx.String(http.StatusInternalServerError, "something went wrong\n")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (h *Handler) GetMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")

	if metricType != string(entity.Gauge) && metricType != string(entity.Counter) {
		ctx.String(http.StatusBadRequest, error2.ErrInvalidType.Error())
		return
	}

	metric, ok := h.uc.GetMetric(metricType, metricName)

	if !ok {
		ctx.String(http.StatusNotFound, error2.ErrNotFound.Error())
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
}

func (h *Handler) GetMetrics(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": h.uc.GetMetrics(),
	})
}
