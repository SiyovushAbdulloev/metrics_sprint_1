package metric

import (
	"bytes"
	"encoding/json"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/metric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupRouterForExtraTests() *gin.Engine {
	db := memory.NewMockDB(make([]entity.Metrics, 0))
	repo := memory.NewMockMetricRepository(db)
	uc := metric.New(repo)
	h := New(uc, nil)

	r := gin.Default()
	r.POST("/value", h.GetMetric)
	r.GET("/value/:type/:name", h.OldGetMetric)
	r.POST("/update/:type/:name/:value", h.OldStoreMetric)
	return r
}

func TestHandler_GetMetric(t *testing.T) {
	r := setupRouterForExtraTests()

	// Pre-store metric
	gaugeValue := 99.9
	preMetric := entity.Metrics{
		ID:    "load",
		MType: entity.Gauge,
		Value: &gaugeValue,
	}
	db := memory.NewMockDB([]entity.Metrics{preMetric})
	repo := memory.NewMockMetricRepository(db)
	uc := metric.New(repo)
	h := New(uc, nil)

	r = gin.Default()
	r.POST("/value", h.GetMetric)

	body, _ := json.Marshal(preMetric)
	req := httptest.NewRequest(http.MethodPost, "/value", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_OldStoreMetric(t *testing.T) {
	r := setupRouterForExtraTests()

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/load/42.42", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_OldGetMetric(t *testing.T) {
	db := memory.NewMockDB([]entity.Metrics{
		{
			ID:    "cpu",
			MType: entity.Gauge,
			Value: ptrFloat64(88.8),
		},
	})
	repo := memory.NewMockMetricRepository(db)
	uc := metric.New(repo)
	h := New(uc, nil)

	r := gin.Default()
	r.GET("/value/:type/:name", h.OldGetMetric)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/cpu", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "88.8")
}

func ptrFloat64(v float64) *float64 { return &v }

func TestHandler_StoreInFileAndRestoreFromFile(t *testing.T) {
	tempFile := "test_metrics_store.txt"
	defer os.Remove(tempFile)

	// Подготовка данных
	db := memory.NewDB([]entity.Metrics{
		{ID: "load_avg", MType: entity.Gauge, Value: ptrFloat64(1.23)},
		{ID: "reqs", MType: entity.Counter, Delta: ptrInt64(100)},
	})
	repo := memory.NewMetricRepository(db)
	uc := metric.New(repo)
	h := New(uc, nil)

	// Сохраняем в файл
	h.StoreInFile(tempFile)

	// Читаем вручную, проверим наличие строк
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	require.Contains(t, string(content), "load_avg")
	require.Contains(t, string(content), "reqs")

	// Новый storage и handler, загружаем туда данные
	emptyRepo := memory.NewMetricRepository(memory.NewDB([]entity.Metrics{}))
	emptyUC := metric.New(emptyRepo)
	emptyHandler := New(emptyUC, nil)

	err = emptyHandler.RestoreFromFile(tempFile)
	require.NoError(t, err)

	metrics, err := emptyUC.GetMetrics()
	require.NoError(t, err)
	require.Len(t, metrics, 2)
	require.Equal(t, "load_avg", metrics[0].ID)
	require.Equal(t, "reqs", metrics[1].ID)
}

func ptrInt64(v int64) *int64 { return &v }
