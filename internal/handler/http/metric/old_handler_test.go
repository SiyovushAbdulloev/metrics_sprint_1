package metric

import (
	"fmt"
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

func TestServer_OldStoreMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type want struct {
		Code        int
		Response    string
		ContentType string
	}
	type Metric struct {
		Type  string
		Name  string
		Value any
	}
	tests := []struct {
		name   string
		metric Metric
		want   want
	}{
		{
			name: "Storing alloc",
			metric: Metric{
				Type:  entity.Gauge,
				Name:  "alloc",
				Value: 1234.32,
			},
			want: want{
				Code:        http.StatusOK,
				Response:    "OK",
				ContentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Storing counter",
			metric: Metric{
				Type:  entity.Counter,
				Name:  "counter",
				Value: 12,
			},
			want: want{
				Code:        http.StatusOK,
				Response:    "OK",
				ContentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Storing unknown type",
			metric: Metric{
				Type:  "counter_2",
				Name:  "anything",
				Value: 12,
			},
			want: want{
				Code:        http.StatusBadRequest,
				Response:    "invalid type",
				ContentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Storing invalid value",
			metric: Metric{
				Type:  entity.Counter,
				Name:  "counter",
				Value: "12asdf",
			},
			want: want{
				Code:        http.StatusBadRequest,
				Response:    "invalid value",
				ContentType: "text/plain; charset=utf-8",
			},
		},
	}

	db := memory.NewMockDB(make([]entity.Metrics, 0))
	metricRepository := memory.NewMockMetricRepository(db)
	metricUC := metric.New(metricRepository)
	hl := New(metricUC, nil)

	router := gin.Default()
	router.POST("/update/:type/:name/:value", hl.OldStoreMetric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := fmt.Sprintf("/update/%s/%s/%v", tt.metric.Type, tt.metric.Name, tt.metric.Value)
			request := httptest.NewRequest(http.MethodPost, addr, nil)
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
