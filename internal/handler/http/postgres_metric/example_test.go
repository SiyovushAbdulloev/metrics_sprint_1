//go:build integration
// +build integration

package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	handler "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/handler/http/metric"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	uc "github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/metric"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_StoreMetric() {
	gin.SetMode(gin.TestMode)
	repo := memory.NewMetricRepository(memory.NewDB(nil))
	usecase := uc.New(repo)
	h := handler.New(usecase, nil)

	r := gin.Default()
	r.POST("/update", h.StoreMetric)

	m := entity.Metrics{
		ID:    "test",
		MType: entity.Gauge,
		Value: ptrFloat64(10.5),
	}
	body, _ := json.Marshal(m)
	req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 200
}

func ExampleHandler_GetMetric() {
	gin.SetMode(gin.TestMode)
	repo := memory.NewMetricRepository(memory.NewDB(nil))
	usecase := uc.New(repo)
	h := handler.New(usecase, nil)

	r := gin.Default()
	r.POST("/update", h.StoreMetric)
	r.POST("/value", h.GetMetric)

	// Сначала сохраняем метрику
	m := entity.Metrics{
		ID:    "test",
		MType: entity.Gauge,
		Value: ptrFloat64(42.42),
	}
	body, _ := json.Marshal(m)
	req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Затем получаем её
	req2 := httptest.NewRequest(http.MethodPost, "/value", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	fmt.Println(w2.Code)
	// Output:
	// 200
}

func ptrFloat64(v float64) *float64 {
	return &v
}
