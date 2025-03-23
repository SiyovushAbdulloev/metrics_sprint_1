package app

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/database/memory"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/models"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/services"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/transport/rest"
)

func Main() {
	db := memory.NewDB(make([]models.Metric, 0))
	metricStorage := memory.NewMetricStorage(db)
	metricService := services.NewMetricService(&metricStorage)
	server := rest.NewServer("localhost", 8080, metricService)

	_, err := server.Run()

	if err != nil {
		panic(err)
	}
}
