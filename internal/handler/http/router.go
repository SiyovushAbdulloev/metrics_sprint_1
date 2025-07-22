package http

import (
	"crypto/rsa"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/config"
	metricHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/middleware"
	checkHandler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"runtime"
)

func DefineMetricRoutes(app *gin.Engine, metricHl *metricHandler.Handler, l logger.Interface, cfg *config.Config) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b))))
	templatesPath := filepath.Join(basePath, "templates", "*.html")

	app.LoadHTMLGlob(templatesPath)

	app.Use(middleware.Logger(l))
	app.Use(middleware.Compress())
	app.Use(middleware.Hash(cfg.App.HashKey))

	app.GET("/", metricHl.GetMetrics)
	app.GET("/value/:type/:name", metricHl.OldGetMetric)
	app.POST("/update/:type/:name/:value", metricHl.OldStoreMetric)
	app.POST("/value/", metricHl.GetMetric)
	app.POST("/update/", metricHl.StoreMetric)
}

func DefinePostgresMetricRoutes(
	app *gin.Engine,
	checkHl *checkHandler.Handler,
	l logger.Interface,
	cfg *config.Config,
	privKey *rsa.PrivateKey,
) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b))))
	templatesPath := filepath.Join(basePath, "templates", "*.html")

	app.LoadHTMLGlob(templatesPath)

	app.Use(middleware.Logger(l))
	app.Use(middleware.Compress())
	app.Use(middleware.Hash(cfg.App.HashKey))

	if privKey != nil {
		app.Use(middleware.DecryptBody(privKey))
	}

	app.GET("/", checkHl.GetMetrics)
	app.GET("/value/:type/:name", checkHl.OldGetMetric)
	app.POST("/update/:type/:name/:value", checkHl.OldStoreMetric)
	app.POST("/value/", checkHl.GetMetric)
	app.POST("/update/", checkHl.StoreMetric)
	app.GET("/ping", checkHl.Check)
	app.POST("/updates", checkHl.UpdateManyMetric)
}
