package rest

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/database"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/services"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"runtime"
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
	server := gin.Default()

	_, b, _, _ := runtime.Caller(0)                                       // Get the current file path
	basePath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b)))) // Go up 3 levels
	templatesPath := filepath.Join(basePath, "templates", "*.html")

	server.LoadHTMLGlob(templatesPath)
	server.GET("/", s.GetMetrics)
	server.GET("/value/:type/:name", s.GetMetric)
	server.POST("/update/:type/:name/:value", s.StoreMetric)

	err := server.Run(s.Addr())

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
