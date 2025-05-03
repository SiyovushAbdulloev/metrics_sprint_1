package metric

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase"
	error2 "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/error"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"io"
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

func (h *Handler) StoreMetric(ctx *gin.Context) {
	var metric entity.Metrics

	body, err := io.ReadAll(ctx.Request.Body)
	fmt.Println("Error before reading all", err)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrSomethingWentWrong.Error(),
		})
		return
	}

	err = easyjson.Unmarshal(body, &metric)
	fmt.Println("Error before unmarshalling all", err)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidValue.Error(),
		})
		return
	}

	if metric.MType != entity.Gauge && metric.MType != entity.Counter {
		fmt.Printf("Error before checking type:")
		ctx.JSON(http.StatusBadRequest, errorResponse{
			Message: error2.ErrInvalidType.Error(),
		})
		return
	}

	_, err = h.uc.StoreMetric(metric)

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

func (h *Handler) StoreInFile(filepath string) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		h.l.Info("Error opening file: %v", "err", err)
		return
	}

	defer file.Close()
	var data bytes.Buffer
	metrics, err := h.uc.GetMetrics()
	if err != nil {
		h.l.Info("Error getting metrics: %v", "err", err)
		return
	}
	for _, metric := range metrics {
		d, err := easyjson.Marshal(metric)
		if err != nil {
			h.l.Info("Error marshalling metric: %v", "err", err)
			break
		}

		d = append(d, '\n')

		data.Write(d)
	}
	file.Write(data.Bytes())
}

func (h *Handler) RestoreFromFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		h.l.Info("Error opening file: %v", "err", err)
		return err
	}

	var metrics []entity.Metrics
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		metric := entity.Metrics{}
		err = easyjson.Unmarshal([]byte(scanner.Text()), &metric)
		if err != nil {
			h.l.Info("Error unmarshalling metric: %v", "err", err)
			return err
		}

		metrics = append(metrics, metric)
	}

	if err = scanner.Err(); err != nil {
		h.l.Info("Error reading file: %v", "err", err)
		return err
	}

	err = h.uc.StoreAll(metrics)

	if err != nil {
		h.l.Info("Error storing metrics")
		return fmt.Errorf("error storing metrics")
	}

	return nil
}
