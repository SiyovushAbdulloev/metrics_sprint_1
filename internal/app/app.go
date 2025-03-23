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
)

func Main(cf *config.Config) {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	db := memory.NewDB(make([]entity.Metric, 0))
	metricRepository := memory.NewMetricRepository(db)
	metricUC := metricUseCase.New(metricRepository)
	metricHl := metricHandler.New(metricUC, l)

	httpServer := httpserver.New(httpserver.WithAddress(cf.Server.Address))
	http.DefineMetricRoutes(httpServer.App, metricHl)

	err = httpServer.Start()

	if err != nil {
		panic(err)
	}

	//var address string
	//addr := os.Getenv("ADDRESS")
	//addrFlag := flag.String("a", "localhost:8080", "The address to listen on for HTTP requests.")
	//
	//flag.Parse()
	//if addr == "" {
	//	address = *addrFlag
	//} else {
	//	address = addr
	//}
	//server := http.InitApp(&metricStorage, address)
	//
	//_, err = server.Run()
	//
	//if err != nil {
	//	panic(err)
	//}
}
