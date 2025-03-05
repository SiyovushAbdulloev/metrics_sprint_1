package app

import (
	"flag"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/database/memory"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/transport/rest"
	"os"
)

func Main() {
	db := memory.NewDB(make([]models.Metric, 0))
	metricStorage := memory.NewMetricStorage(db)
	var address string
	addr := os.Getenv("ADDRESS")
	addrFlag := flag.String("a", "localhost:8080", "The address to listen on for HTTP requests.")

	flag.Parse()
	if addr == "" {
		address = *addrFlag
	} else {
		address = addr
	}
	server := rest.InitApp(&metricStorage, address)

	_, err := server.Run()

	if err != nil {
		panic(err)
	}
}
