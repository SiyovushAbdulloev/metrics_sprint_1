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
	"log"
	"net/http"
	"os"
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

	log.SetOutput(os.Stdout)

	body, err := io.ReadAll(ctx.Request.Body)
	log.Println("Body", string(body))
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

	log.Println("Metric", metric)

	if metric.MType != entity.Gauge && metric.MType != entity.Counter {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidType.Error(),
		})
		return
	}

	_, err = h.uc.StoreMetric(metric)

	log.Println("Err:", err)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, metric)
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

	m, err := h.uc.GetMetric(metric)

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
