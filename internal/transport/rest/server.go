package rest

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/database"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/services"
	"net/http"
)

type Server struct {
	Host    string
	Port    int16
	Service services.MetricService
}

func InitApp(metricStorage database.MetricStorage) *Server {
	metricService := services.NewMetricService(metricStorage)
	return NewServer("localhost", 8080, metricService)
}

func NewServer(host string, port int16, service services.MetricService) *Server {
	return &Server{
		Host:    host,
		Port:    port,
		Service: service,
	}
}

func (s *Server) Run() (bool, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/{type}/{name}/{value}", s.StoreMetric)

	err := http.ListenAndServe(s.Addr(), mux)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
