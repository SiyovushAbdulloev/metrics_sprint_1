package app

import (
	"database/sql"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/config"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	handler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http"
	metricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/metric"
	postgresMetricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	postgresRepo "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/postgres"
	metricUseCase "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/metric"
	postgresMetricUseCase "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/crypto"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/httpserver"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func Main(cf *config.Config) {
	binding.Validator = nil
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
	if cf.Database.DSN != "" {
		dbMigration, err2 := sql.Open("postgres", cf.Database.DSN)
		if err2 != nil {
			l.Info("failed to open postgres db:", "err", err2)
			panic(err2)
		}
		defer dbMigration.Close()

		if err2 = goose.SetDialect("postgres"); err2 != nil {
			l.Info("failed to set goose dialect:", "err", err)
			panic(err2)
		}

		if err2 = goose.Up(dbMigration, "./migrations"); err2 != nil {
			l.Info("failed to migrate goose migrations:", "err", err2)
			panic(err2)
		}
	}

	metricRepository := memory.NewMetricRepository(db)
	pgRepo := postgresRepo.NewMetricRepository(postgresDB)
	metricUC := metricUseCase.New(metricRepository)
	postgresUC := postgresMetricUseCase.New(pgRepo)
	metricHl := metricHandler.New(metricUC, l)
	postgresHl := postgresMetricHandler.New(postgresUC, l)

	httpServer := httpserver.New(httpserver.WithAddress(cf.Server.Address))

	if cf.Database.DSN != "" {
		key, err := crypto.LoadPrivateKey(cf.App.CryptoKeyPath)
		if err == nil {
			handler.DefinePostgresMetricRoutes(httpServer.App, postgresHl, l, cf, key)
		} else {
			handler.DefinePostgresMetricRoutes(httpServer.App, postgresHl, l, cf, nil)
		}
	} else {
		handler.DefineMetricRoutes(httpServer.App, metricHl, l, cf)
	}

	if cf.App.Restore {
		err = metricHl.RestoreFromFile(cf.App.Filepath)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
