package http

import (
	metricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/metric"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"runtime"
)

func DefineMetricRoutes(app *gin.Engine, metricHl *metricHandler.Handler) {
	_, b, _, _ := runtime.Caller(0)                                       // Get the current file path
	basePath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b)))) // Go up 3 levels
	templatesPath := filepath.Join(basePath, "templates", "*.html")

	app.LoadHTMLGlob(templatesPath)

	app.GET("/", metricHl.GetMetrics)
	app.GET("/value/:type/:name", metricHl.GetMetric)
	app.POST("/update/:type/:name/:value", metricHl.StoreMetric)
}

//type Server struct {
//	Address string
//	Service services.MetricService
//}

//func InitApp(metricStorage repository.MetricStorage, address string) *Server {
//	metricService := services.NewMetricService(metricStorage)
//	return NewServer(address, metricService)
//}
//
//func NewServer(address string, service services.MetricService) *Server {
//	return &Server{
//		Address: address,
//		Service: service,
//	}
//}

//func (s *Server) Run() (bool, error) {
//	server := gin.Default()
//
//	_, b, _, _ := runtime.Caller(0)                                       // Get the current file path
//	basePath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b)))) // Go up 3 levels
//	templatesPath := filepath.Join(basePath, "templates", "*.html")
//
//	server.LoadHTMLGlob(templatesPath)
//	server.GET("/", s.GetMetrics)
//	server.GET("/value/:type/:name", s.GetMetric)
//	server.POST("/update/:type/:name/:value", s.StoreMetric)
//
//	err := server.Run(s.Address)
//
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
