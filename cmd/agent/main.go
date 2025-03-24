package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Metric struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Value float64 `json:"value"`
	Delta int64   `json:"delta"`
}

type Metrics struct {
	data []Metric
}

type Config struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func collectMetrics(m *Metrics) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	var counter int64 = 1

	if length := len(m.data); length > 0 {
		counter = m.data[len(m.data)-1].Delta + 1
	}

	data := []Metric{
		{
			ID:    "alloc",
			MType: "gauge",
			Value: float64(rtm.Alloc),
		},
		{
			ID:    "buck_hash_sys",
			MType: "gauge",
			Value: float64(rtm.BuckHashSys),
		},
		{
			ID:    "frees",
			MType: "gauge",
			Value: float64(rtm.Frees),
		},
		{
			ID:    "gccpu_fraction",
			MType: "gauge",
			Value: float64(rtm.GCCPUFraction),
		},
		{
			ID:    "heap_alloc",
			MType: "gauge",
			Value: float64(rtm.HeapAlloc),
		},
		{
			ID:    "heap_idle",
			MType: "gauge",
			Value: float64(rtm.HeapIdle),
		},
		{
			ID:    "heap_inuse",
			MType: "gauge",
			Value: float64(rtm.HeapInuse),
		},
		{
			ID:    "heap_objects",
			MType: "gauge",
			Value: float64(rtm.HeapObjects),
		},
		{
			ID:    "heap_released",
			MType: "gauge",
			Value: float64(rtm.HeapReleased),
		},
		{
			ID:    "heap_sys",
			MType: "gauge",
			Value: float64(rtm.HeapSys),
		},
		{
			ID:    "last_gc",
			MType: "gauge",
			Value: float64(rtm.LastGC),
		},
		{
			ID:    "lookups",
			MType: "gauge",
			Value: float64(rtm.Lookups),
		},
		{
			ID:    "mcache_inuse",
			MType: "gauge",
			Value: float64(rtm.MCacheInuse),
		},
		{
			ID:    "m_cache_sys",
			MType: "gauge",
			Value: float64(rtm.MCacheSys),
		},
		{
			ID:    "mspan_inuse",
			MType: "gauge",
			Value: float64(rtm.MSpanInuse),
		},
		{
			ID:    "mspan_sys",
			MType: "gauge",
			Value: float64(rtm.MSpanSys),
		},
		{
			ID:    "mallocs",
			MType: "gauge",
			Value: float64(rtm.Mallocs),
		},
		{
			ID:    "next_gc",
			MType: "gauge",
			Value: float64(rtm.NextGC),
		},
		{
			ID:    "num_forced_gc",
			MType: "gauge",
			Value: float64(rtm.NumForcedGC),
		},
		{
			ID:    "num_gc",
			MType: "gauge",
			Value: float64(rtm.NumGC),
		},
		{
			ID:    "other_sys",
			MType: "gauge",
			Value: float64(rtm.OtherSys),
		},
		{
			ID:    "pause_totalns",
			MType: "gauge",
			Value: float64(rtm.PauseTotalNs),
		},
		{
			ID:    "stack_inuse",
			MType: "gauge",
			Value: float64(rtm.StackInuse),
		},
		{
			ID:    "stack_sys",
			MType: "gauge",
			Value: float64(rtm.StackSys),
		},
		{
			ID:    "sys",
			MType: "gauge",
			Value: float64(rtm.Sys),
		},
		{
			ID:    "total_alloc",
			MType: "gauge",
			Value: float64(rtm.TotalAlloc),
		},
		{
			ID:    "random_value",
			MType: "gauge",
			Value: float64(1),
		},
		{
			ID:    "poll_count",
			MType: "counter",
			Delta: counter,
		},
	}

	m.data = data
}

func sendMetrics(client http.Client, m Metrics, address string) {
	for _, metric := range m.data {
		data, err := json.Marshal(metric)
		if err != nil {
			log.Printf("Error marshaling metric: %v", err)
			return
		}

		var buffer bytes.Buffer
		gzipW, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
		if err != nil {
			fmt.Printf("Error creating gzip writer: %v", err)
			return
		}

		_, err = gzipW.Write(data)
		if err != nil {
			fmt.Printf("Error writing metric: %v", err)
			return
		}

		gzipW.Close()

		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/%s", address, "update"), &buffer)
		if err != nil {
			fmt.Printf("Error creating request: %v", err)
			return
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error posting metric: %v", err)
		}

		res.Body.Close()
	}
}

func getConfig() Config {
	var address string
	var reportInterval int
	var pollInterval int

	addr := os.Getenv("ADDRESS")
	reportInt := os.Getenv("REPORT_INTERVAL")
	pollInt := os.Getenv("POLL_INTERVAL")

	addrFlag := flag.String("a", "localhost:8080", "The address to send HTTP requests.")
	reportIntFlag := flag.Int("r", 10, "The interval in seconds between metric reporting. (in seconds)")
	pollIntFlag := flag.Int("p", 2, "The interval in seconds between metric polling. (in seconds)")
	flag.Parse()

	if addr == "" {
		address = *addrFlag
	} else {
		address = addr
	}

	if reportInt == "" {
		reportInterval = *reportIntFlag
	} else {
		value, err := strconv.Atoi(reportInt)
		if err != nil {
			panic(err)
		}
		reportInterval = value
	}

	if pollInt == "" {
		pollInterval = *pollIntFlag
	} else {
		value, err := strconv.Atoi(pollInt)
		if err != nil {
			panic(err)
		}
		pollInterval = value
	}

	return Config{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}

func main() {
	config := getConfig()
	client := http.Client{}
	m := Metrics{}
	collectTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer collectTicker.Stop()
	defer sendTicker.Stop()

	go func() {
		for {
			select {
			case <-collectTicker.C:
				collectMetrics(&m)
			case <-sendTicker.C:
				sendMetrics(client, m, config.Address)
			}
		}
	}()

	select {}
}
