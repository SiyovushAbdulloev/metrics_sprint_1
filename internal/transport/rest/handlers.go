package rest

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) StoreMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")
	metricValue := ctx.Param("value")

	if metricType != string(models.Gauge) && metricType != string(models.Counter) {
		ctx.String(http.StatusBadRequest, errInvalidType.Error())
		return
	}

	value, ok := s.validValue(metricType, metricValue)
	if !ok {
		ctx.String(http.StatusBadRequest, errInvalidValue.Error())
		return
	}

	added := s.Service.StoreMetric(models.Metric{
		Name:  metricName,
		Value: value,
		Type:  models.MetricType(metricType),
	})

	if !added {
		ctx.String(http.StatusInternalServerError, "something went wrong\n")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (s *Server) GetMetric(ctx *gin.Context) {
	metricType := ctx.Param("type")
	metricName := ctx.Param("name")

	if metricType != string(models.Gauge) && metricType != string(models.Counter) {
		ctx.String(http.StatusBadRequest, errInvalidType.Error())
		return
	}

	metric, ok := s.Service.GetMetric(metricType, metricName)

	if !ok {
		ctx.String(http.StatusNotFound, errNotFound.Error())
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
}

func (s *Server) GetMetrics(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": s.Service.GetMetrics(),
	})
}
