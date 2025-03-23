package rest

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
	"net/http"
)

func (s *Server) StoreMetric(w http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue("type")
	metricName := req.PathValue("name")
	metricValue := req.PathValue("value")

	fmt.Printf("TYPE: %s\n", metricType)

	if metricType != string(models.Gauge) && metricType != string(models.Counter) {
		http.Error(w, errInvalidType.Error(), http.StatusBadRequest)
		return
	}

	value, ok := s.validValue(metricType, metricValue)
	if !ok {
		http.Error(w, errInvalidValue.Error(), http.StatusBadRequest)
		return
	}

	added := s.Service.StoreMetric(models.Metric{
		Name:  metricName,
		Value: value,
		Type:  models.MetricType(metricType),
	})

	if !added {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
