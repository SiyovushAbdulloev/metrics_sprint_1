package main

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/config"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/app"
)

func main() {
	cf, err := config.New()
	if err != nil {
		panic(err)
	}

	app.Main(cf)
}
