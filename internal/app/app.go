package app

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/config"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http"
	metricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	metricUseCase "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/httpserver"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"log"
	"os"
)

func Main(cf *config.Config) {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	db := memory.NewDB(make([]entity.Metrics, 0))
	metricRepository := memory.NewMetricRepository(db)
	metricUC := metricUseCase.New(metricRepository)
	metricHl := metricHandler.New(metricUC, l)

	httpServer := httpserver.New(httpserver.WithAddress(cf.Server.Address))
	http.DefineMetricRoutes(httpServer.App, metricHl, l)

	log.SetOutput(os.Stdout)
	log.Println("Starting server on " + cf.Server.Address)
	err = httpServer.Start()
	log.Println("Starting failed with error: ", err)

	if err != nil {
		panic(err)
	}
}
