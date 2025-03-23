package metric

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase"
	error2 "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"io"
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
	var metric entity.Metrics

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	err = easyjson.Unmarshal(body, &metric)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidValue.Error(),
		})
		return
	}

	if metric.MType != entity.Gauge && metric.MType != entity.Counter {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidType.Error(),
		})
		return
	}

	added := h.uc.StoreMetric(metric)

	if added.ID == "" {
		ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, metric)
}

func (h *Handler) GetMetric(ctx *gin.Context) {
	var metric entity.Metrics
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	err = easyjson.Unmarshal(body, &metric)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidValue.Error(),
		})
	}

	if metric.MType != entity.Gauge && metric.MType != entity.Counter {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidType.Error(),
		})
		return
	}

	m, ok := h.uc.GetMetric(metric)

	if !ok {
		ctx.JSON(http.StatusNotFound, errorResponse{
			Message: error2.ErrNotFound.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, m)
}

func (h *Handler) GetMetrics(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": h.uc.GetMetrics(),
	})
}
