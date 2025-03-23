package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

type Metric struct {
	Name  string
	Type  string
	Value any
}

type Metrics struct {
	data []Metric
}

func collectMetrics(m *Metrics) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	var counter int64 = 1

	if length := len(m.data); length > 0 {
		counter = m.data[len(m.data)-1].Value.(int64) + 1
	}

	data := []Metric{
		{
			Name:  "alloc",
			Type:  "gauge",
			Value: float64(rtm.Alloc),
		},
		{
			Name:  "buck_hash_sys",
			Type:  "gauge",
			Value: float64(rtm.BuckHashSys),
		},
		{
			Name:  "frees",
			Type:  "gauge",
			Value: float64(rtm.Frees),
		},
		{
			Name:  "gccpu_fraction",
			Type:  "gauge",
			Value: float64(rtm.GCCPUFraction),
		},
		{
			Name:  "heap_alloc",
			Type:  "gauge",
			Value: float64(rtm.HeapAlloc),
		},
		{
			Name:  "heap_idle",
			Type:  "gauge",
			Value: float64(rtm.HeapIdle),
		},
		{
			Name:  "heap_inuse",
			Type:  "gauge",
			Value: float64(rtm.HeapInuse),
		},
		{
			Name:  "heap_objects",
			Type:  "gauge",
			Value: float64(rtm.HeapObjects),
		},
		{
			Name:  "heap_released",
			Type:  "gauge",
			Value: float64(rtm.HeapReleased),
		},
		{
			Name:  "heap_sys",
			Type:  "gauge",
			Value: float64(rtm.HeapSys),
		},
		{
			Name:  "last_gc",
			Type:  "gauge",
			Value: float64(rtm.LastGC),
		},
		{
			Name:  "lookups",
			Type:  "gauge",
			Value: float64(rtm.Lookups),
		},
		{
			Name:  "mcache_inuse",
			Type:  "gauge",
			Value: float64(rtm.MCacheInuse),
		},
		{
			Name:  "m_cache_sys",
			Type:  "gauge",
			Value: float64(rtm.MCacheSys),
		},
		{
			Name:  "mspan_inuse",
			Type:  "gauge",
			Value: float64(rtm.MSpanInuse),
		},
		{
			Name:  "mspan_sys",
			Type:  "gauge",
			Value: float64(rtm.MSpanSys),
		},
		{
			Name:  "mallocs",
			Type:  "gauge",
			Value: float64(rtm.Mallocs),
		},
		{
			Name:  "next_gc",
			Type:  "gauge",
			Value: float64(rtm.NextGC),
		},
		{
			Name:  "num_forced_gc",
			Type:  "gauge",
			Value: float64(rtm.NumForcedGC),
		},
		{
			Name:  "num_gc",
			Type:  "gauge",
			Value: float64(rtm.NumGC),
		},
		{
			Name:  "other_sys",
			Type:  "gauge",
			Value: float64(rtm.OtherSys),
		},
		{
			Name:  "pause_totalns",
			Type:  "gauge",
			Value: float64(rtm.PauseTotalNs),
		},
		{
			Name:  "stack_inuse",
			Type:  "gauge",
			Value: float64(rtm.StackInuse),
		},
		{
			Name:  "stack_sys",
			Type:  "gauge",
			Value: float64(rtm.StackSys),
		},
		{
			Name:  "sys",
			Type:  "gauge",
			Value: float64(rtm.Sys),
		},
		{
			Name:  "total_alloc",
			Type:  "gauge",
			Value: float64(rtm.TotalAlloc),
		},
		{
			Name:  "random_value",
			Type:  "gauge",
			Value: float64(1),
		},
		{
			Name:  "poll_count",
			Type:  "counter",
			Value: int64(counter),
		},
	}

	m.data = data
}

func sendMetrics(client http.Client, m Metrics) {
	for _, metric := range m.data {
		addr := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metric.Type, metric.Name, metric.Value)
		res, err := client.Post(addr, "text/plain", nil)
		if err != nil {
			log.Printf("Error posting metric: %v", err)
		}

		res.Body.Close()
	}
}

func main() {
	client := http.Client{}
	m := Metrics{}
	collectTicker := time.NewTicker(2 * time.Second)
	sendTicker := time.NewTicker(10 * time.Second)
	defer collectTicker.Stop()
	defer sendTicker.Stop()

	go func() {
		for {
			select {
			case <-collectTicker.C:
				collectMetrics(&m)
			case <-sendTicker.C:
				sendMetrics(client, m)
			}
		}
	}()

	select {}
}
