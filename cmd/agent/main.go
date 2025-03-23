package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Metric struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Value float64 `json:"value,omitempty"`
	Delta int64   `json:"delta,omitempty"`
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
			ID:    "Alloc",
			MType: "gauge",
			Value: float64(rtm.Alloc),
		},
		{
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: float64(rtm.BuckHashSys),
		},
		{
			ID:    "Frees",
			MType: "gauge",
			Value: float64(rtm.Frees),
		},
		{
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: float64(rtm.GCCPUFraction),
		},
		{
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: float64(rtm.HeapAlloc),
		},
		{
			ID:    "HeapInuse",
			MType: "gauge",
			Value: float64(rtm.HeapInuse),
		},
		{
			ID:    "HeapObjects",
			MType: "gauge",
			Value: float64(rtm.HeapObjects),
		},
		{
			ID:    "HeapReleased",
			MType: "gauge",
			Value: float64(rtm.HeapReleased),
		},
		{
			ID:    "HeapSys",
			MType: "gauge",
			Value: float64(rtm.HeapSys),
		},
		{
			ID:    "LastGC",
			MType: "gauge",
			Value: float64(rtm.LastGC),
		},
		{
			ID:    "Lookups",
			MType: "gauge",
			Value: float64(rtm.Lookups),
		},
		{
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: float64(rtm.MCacheInuse),
		},
		{
			ID:    "MCacheSys",
			MType: "gauge",
			Value: float64(rtm.MCacheSys),
		},
		{
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: float64(rtm.MSpanInuse),
		},
		{
			ID:    "MSpanSys",
			MType: "gauge",
			Value: float64(rtm.MSpanSys),
		},
		{
			ID:    "Mallocs",
			MType: "gauge",
			Value: float64(rtm.Mallocs),
		},
		{
			ID:    "NextGC",
			MType: "gauge",
			Value: float64(rtm.NextGC),
		},
		{
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: float64(rtm.NumForcedGC),
		},
		{
			ID:    "NumGC",
			MType: "gauge",
			Value: float64(rtm.NumGC),
		},
		{
			ID:    "OtherSys",
			MType: "gauge",
			Value: float64(rtm.OtherSys),
		},
		{
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: float64(rtm.PauseTotalNs),
		},
		{
			ID:    "StackInuse",
			MType: "gauge",
			Value: float64(rtm.StackInuse),
		},
		{
			ID:    "StackSys",
			MType: "gauge",
			Value: float64(rtm.StackSys),
		},
		{
			ID:    "Sys",
			MType: "gauge",
			Value: float64(rtm.Sys),
		},
		{
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: float64(rtm.TotalAlloc),
		},
		{
			ID:    "GCSys",
			MType: "gauge",
			Value: float64(rtm.GCSys),
		},
		{
			ID:    "HeapIdle",
			MType: "gauge",
			Value: float64(rtm.HeapIdle),
		},
		{
			ID:    "RandomValue",
			MType: "gauge",
			Value: rand.Float64(),
		},
		{
			ID:    "PollCount",
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
		body := bytes.NewBuffer(data)

		res, err := client.Post(fmt.Sprintf("http://%s/update/", address), "application/json", body)
		if err != nil {
			log.Printf("Error posting metric: %v", err)
		}

		if err == nil {
			res.Body.Close()
		}
	}
}

func getVars() Config {
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
	config := getVars()
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
