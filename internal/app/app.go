package app

import (
	"flag"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/database/memory"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/transport/rest"
)

func Main() {
	db := memory.NewDB(make([]models.Metric, 0))
	metricStorage := memory.NewMetricStorage(db)
	address := flag.String("a", "localhost:8080", "The address to listen on for HTTP requests.")

	flag.Parse()
	server := rest.InitApp(&metricStorage, *address)

	_, err := server.Run()

	if err != nil {
		panic(err)
	}
}
