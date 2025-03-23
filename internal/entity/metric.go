package entity

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func New(id string, metricType string, delta *int64, value *float64) *Metrics {
	return &Metrics{
		ID:    id,
		MType: metricType,
		Delta: delta,
		Value: value,
	}
}
