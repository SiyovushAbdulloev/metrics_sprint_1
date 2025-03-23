package metric

import (
	"bytes"
	"encoding/json"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/entity"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/repository/memory"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/internal/usecase/metric"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_StoreMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type want struct {
		Code        int
		Response    any
		ContentType string
	}
	type Metric struct {
		ID    string  `json:"id"`
		MType string  `json:"type"`
		Delta int64   `json:"delta"`
		Value float64 `json:"value"`
	}
	tests := []struct {
		name   string
		metric Metric
		want   want
	}{
		{
			name: "Storing alloc",
			metric: Metric{
				MType: entity.Gauge,
				ID:    "alloc",
				Value: 1234.32,
			},
			want: want{
				Code:        http.StatusOK,
				Response:    "{\"id\":\"alloc\",\"type\":\"gauge\",\"delta\":0,\"value\":1234.32}",
				ContentType: "application/json; charset=utf-8",
			},
		},
		{
			name: "Storing counter",
			metric: Metric{
				MType: entity.Counter,
				ID:    "counter",
				Delta: 12,
			},
			want: want{
				Code:        http.StatusOK,
				Response:    "{\"id\":\"counter\",\"type\":\"counter\",\"delta\":12,\"value\":0}",
				ContentType: "application/json; charset=utf-8",
			},
		},
		{
			name: "Storing unknown type",
			metric: Metric{
				MType: "counter_2",
				ID:    "anything",
				Value: 12,
			},
			want: want{
				Code:        http.StatusBadRequest,
				Response:    "{\"message\":\"invalid type\"}",
				ContentType: "application/json; charset=utf-8",
			},
		},
	}

	db := memory.NewMockDB(make([]entity.Metrics, 0))
	metricRepository := memory.NewMockMetricRepository(db)
	metricUC := metric.New(metricRepository)
	hl := New(metricUC, nil)

	router := gin.Default()
	router.POST("/update", hl.StoreMetric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := json.Marshal(tt.metric)
			body := bytes.NewBuffer(data)
			request := httptest.NewRequest(http.MethodPost, "/update", body)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.Code, res.StatusCode)

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.Response, string(resBody))
			assert.Equal(t, tt.want.ContentType, res.Header.Get("Content-Type"))
		})
	}
}
