package metric

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase"
	error2 "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"slices"
	"strconv"
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

func (h *Handler) Check(ctx *gin.Context) {
	err := h.uc.Check()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
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

	h.l.Info("Storing metric", "metric", metric)
	metrics, _ := h.uc.GetMetrics()
	h.l.Info("AllMetrics", "allMetrics", metrics)
	_, err = h.uc.StoreMetric(metric)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, metric)
}

func (h *Handler) UpdateManyMetric(ctx *gin.Context) {
	var metricsList entity.MetricsList

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidValue.Error(),
		})
		return
	}

	err = easyjson.Unmarshal(body, &metricsList)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidValue.Error(),
		})
		return
	}

	if len(metricsList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
		return
	}

	metricsIds := make([]string, 0, len(metricsList))
	newMetrics := make([]entity.Metrics, 0, len(metricsList))
	counters := make(map[string]int64)
	for _, m := range metricsList {
		if m.MType == entity.Counter {
			prev, ok := counters[m.ID]
			var delta int64
			if ok {
				delta = prev
			}
			counters[m.ID] = delta + *m.Delta
		} else {
			if !slices.Contains(metricsIds, m.ID) {
				newMetrics = append(newMetrics, m)
				metricsIds = append(metricsIds, m.ID)
			} else {
				for k, mt := range newMetrics {
					if mt.ID == m.ID {
						newMetrics[k] = m
					}
				}
			}
		}
	}

	for id, delta := range counters {
		newMetrics = append(newMetrics, entity.Metrics{
			MType: entity.Counter,
			ID:    id,
			Delta: &delta,
		})
	}

	err = h.uc.UpdateAll(newMetrics)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Ok",
	})
}

func (h *Handler) OldStoreMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")
	metricValue := ctx.Param("value")

	if metricType != entity.Gauge && metricType != entity.Counter {
		ctx.String(http.StatusBadRequest, error2.ErrInvalidType.Error())
		return
	}

	_, ok := h.validValue(metricType, metricValue)
	if !ok {
		ctx.String(http.StatusBadRequest, error2.ErrInvalidValue.Error())
		return
	}

	var metric entity.Metrics
	if metricType == entity.Gauge {
		value, _ := strconv.ParseFloat(metricValue, 64)
		metric = entity.Metrics{
			ID:    metricName,
			Value: &value,
			MType: metricType,
		}
	} else if metricType == entity.Counter {
		delta, _ := strconv.ParseInt(metricValue, 10, 64)
		metric = entity.Metrics{
			ID:    metricName,
			Delta: &delta,
			MType: metricType,
		}
	}

	_, err := h.uc.StoreMetric(metric)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "something went wrong\n")
		return
	}

	ctx.String(http.StatusOK, "OK")
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

	h.l.Info("GetMetric", "metric", metric)
	m, err := h.uc.GetMetric(metric)

	metrics, _ := h.uc.GetMetrics()
	h.l.Info("AllMetrics", "allMetrics", metrics)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse{
			Message: error2.ErrNotFound.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, m)
}

func (h *Handler) OldGetMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")

	if metricType != entity.Gauge && metricType != entity.Counter {
		ctx.String(http.StatusBadRequest, error2.ErrInvalidType.Error())
		return
	}

	metric, err := h.uc.GetMetric(entity.Metrics{
		MType: metricType,
		ID:    metricName,
	})

	if err != nil {
		ctx.String(http.StatusNotFound, error2.ErrNotFound.Error())
		return
	}

	var value any
	if metric.MType == entity.Gauge {
		value = *metric.Value
	} else if metric.MType == entity.Counter {
		value = *metric.Delta
	}

	ctx.String(http.StatusOK, fmt.Sprintf("%v", value))
}

func (h *Handler) GetMetrics(ctx *gin.Context) {
	metrics, err := h.uc.GetMetrics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": metrics,
	})
}
