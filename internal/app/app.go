package app

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/config"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http"
	metricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/metric"
	postgresMetricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	postgresRepo "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/postgres"
	metricUseCase "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/metric"
	postgresMetricUseCase "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/httpserver"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
	"time"
)

func Main(cf *config.Config) {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	db := memory.NewDB(make([]entity.Metrics, 0))
	postgresDB, err := postgres.New(cf.Database.DSN)
	if err != nil {
		l.Info("Error connecting to database", "err", err)
		panic(err)
	}
	metricRepository := memory.NewMetricRepository(db)
	pgRepo := postgresRepo.NewMetricRepository(postgresDB)
	metricUC := metricUseCase.New(metricRepository)
	postgresUC := postgresMetricUseCase.New(pgRepo)
	metricHl := metricHandler.New(metricUC, l)
	postgresHl := postgresMetricHandler.New(postgresUC, l)

	httpServer := httpserver.New(httpserver.WithAddress(cf.Server.Address))

	if cf.Database.DSN != "" {
		http.DefinePostgresMetricRoutes(httpServer.App, postgresHl, l)
	} else {
		http.DefineMetricRoutes(httpServer.App, metricHl, l)
	}

	if cf.App.Restore {
		err = metricHl.RestoreFromFile(cf.App.Filepath)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		err = httpServer.Start()

		if err != nil {
			panic(err)
		}
	}()

	storeTicker := time.NewTicker(time.Second * time.Duration(cf.App.StoreInterval))
	defer storeTicker.Stop()

	go func() {
		for range storeTicker.C {
			metricHl.StoreInFile(cf.App.Filepath)
		}
	}()

	select {}
}
