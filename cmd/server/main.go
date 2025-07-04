package main

import (
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/config"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/app"
)

var buildVersion string = "1.0.0"
var buildDate string = "2025-07-04"
var buildCommit string = "HEAD"

func main() {
	fmt.Printf("Build version: %s (или \"N/A\" при отсутствии значения) \n", buildVersion)
	fmt.Printf("Build date: %s (или \"N/A\" при отсутствии значения) \n", buildDate)
	fmt.Printf("Build commit: %s (или \"N/A\" при отсутствии значения) \n", buildCommit)
	cf, err := config.New()
	if err != nil {
		panic(err)
	}

	app.Main(cf)
}
