// Package entity содержит определения сущностей, используемых в сервисе метрик.
package entity

const (
	// Gauge — тип метрики, представляющий значение с плавающей запятой.
	Gauge string = "gauge"
	// Counter — тип метрики, представляющий целочисленный счётчик.
	Counter string = "counter"
)

// Metrics описывает метрику, которая может быть типа gauge или counter.
//
// ID — имя метрики (уникальный идентификатор).
// MType — тип метрики: "gauge" или "counter".
// Delta — значение счётчика (используется, если MType == "counter").
// Value — значение gauge (используется, если MType == "gauge").
//
//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// MetricsList представляет собой список метрик.
//
//easyjson:json
type MetricsList []Metrics
