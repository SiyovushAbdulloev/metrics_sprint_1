package metric_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/postgres"
	metricUC "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/postgres_metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	pgpkg "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	dsn     = flag.String("dsn", "postgres://postgres:password@localhost:5432/metrics", "PostgreSQL DSN")
	db      *pgpkg.Postgres
	cleaned = false
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *dsn == "" {
		panic("DSN must be provided with --dsn")
	}
	var err error
	db, err = pgpkg.New(*dsn)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func cleanupMetricsTable(t *testing.T) {
	if db == nil {
		t.Fatal("Postgres connection is not initialized")
	}
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM metrics")
	require.NoError(t, err)
}

func setupHandler(t *testing.T) *gin.Engine {
	if !cleaned {
		cleanupMetricsTable(t)
		cleaned = true
	}

	repo := postgres.NewMetricRepository(db)
	uc := metricUC.New(repo)
	log, _ := logger.New()
	h := metric.New(uc, log)

	r := gin.Default()
	r.POST("/update", h.StoreMetric)
	r.GET("/ping", h.Check)
	r.POST("/updates", h.UpdateManyMetric)
	r.POST("/value", h.GetMetric)
	// r.GET("/values", h.GetMetrics) // Удалён, чтобы избежать паники из-за отсутствующего шаблона

	return r
}

func TestHandler_StoreMetric(t *testing.T) {
	router := setupHandler(t)

	tests := []struct {
		name     string
		input    entity.Metrics
		wantCode int
	}{
		{
			name: "Store Gauge",
			input: entity.Metrics{
				ID:    "test_gauge",
				MType: entity.Gauge,
				Value: ptrFloat64(100.5),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Store Counter",
			input: entity.Metrics{
				ID:    "test_counter",
				MType: entity.Counter,
				Delta: ptrInt64(42),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Invalid Type",
			input: entity.Metrics{
				ID:    "bad",
				MType: "wrong",
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tt.wantCode, rec.Code)
		})
	}
}

func TestHandler_Check(t *testing.T) {
	router := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UpdateManyMetric(t *testing.T) {
	router := setupHandler(t)

	list := entity.MetricsList{
		{
			ID:    "cpu",
			MType: entity.Gauge,
			Value: ptrFloat64(90.5),
		},
		{
			ID:    "ops",
			MType: entity.Counter,
			Delta: ptrInt64(12),
		},
	}

	body, _ := json.Marshal(list)
	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func ptrFloat64(v float64) *float64 { return &v }
func ptrInt64(v int64) *int64       { return &v }
