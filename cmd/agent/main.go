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
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
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
		counter = *(m.data[len(m.data)-1].Delta) + 1
	}

	alloc := float64(rtm.Alloc)
	buckHashSys := float64(rtm.BuckHashSys)
	frees := float64(rtm.Frees)
	gccpuFraction := float64(rtm.GCCPUFraction)
	heapAlloc := float64(rtm.HeapAlloc)
	heapInUse := float64(rtm.HeapInuse)
	heapObjects := float64(rtm.HeapObjects)
	heapReleased := float64(rtm.HeapReleased)
	heapSys := float64(rtm.HeapSys)
	lastGc := float64(rtm.LastGC)
	lookups := float64(rtm.Lookups)
	mCacheInUse := float64(rtm.MCacheInuse)
	mCacheSys := float64(rtm.MCacheSys)
	mSpanInuse := float64(rtm.MSpanInuse)
	mSpanSys := float64(rtm.MSpanSys)
	mallocs := float64(rtm.Mallocs)
	nextGC := float64(rtm.NextGC)
	numForcedGC := float64(rtm.NumForcedGC)
	numGC := float64(rtm.NumGC)
	otherSys := float64(rtm.OtherSys)
	pauseTotalNs := float64(rtm.PauseTotalNs)
	stackInUse := float64(rtm.StackInuse)
	stackSys := float64(rtm.StackSys)
	sys := float64(rtm.Sys)
	totalAlloc := float64(rtm.TotalAlloc)
	gcSys := float64(rtm.GCSys)
	heapIdle := float64(rtm.HeapIdle)
	randomValue := rand.Float64()

	data := []Metric{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &alloc,
		},
		{
			ID:    "BuckHashSys",
			MType: "gauge",
			Value: &buckHashSys,
		},
		{
			ID:    "Frees",
			MType: "gauge",
			Value: &frees,
		},
		{
			ID:    "GCCPUFraction",
			MType: "gauge",
			Value: &gccpuFraction,
		},
		{
			ID:    "HeapAlloc",
			MType: "gauge",
			Value: &heapAlloc,
		},
		{
			ID:    "HeapInuse",
			MType: "gauge",
			Value: &heapInUse,
		},
		{
			ID:    "HeapObjects",
			MType: "gauge",
			Value: &heapObjects,
		},
		{
			ID:    "HeapReleased",
			MType: "gauge",
			Value: &heapReleased,
		},
		{
			ID:    "HeapSys",
			MType: "gauge",
			Value: &heapSys,
		},
		{
			ID:    "LastGC",
			MType: "gauge",
			Value: &lastGc,
		},
		{
			ID:    "Lookups",
			MType: "gauge",
			Value: &lookups,
		},
		{
			ID:    "MCacheInuse",
			MType: "gauge",
			Value: &mCacheInUse,
		},
		{
			ID:    "MCacheSys",
			MType: "gauge",
			Value: &mCacheSys,
		},
		{
			ID:    "MSpanInuse",
			MType: "gauge",
			Value: &mSpanInuse,
		},
		{
			ID:    "MSpanSys",
			MType: "gauge",
			Value: &mSpanSys,
		},
		{
			ID:    "Mallocs",
			MType: "gauge",
			Value: &mallocs,
		},
		{
			ID:    "NextGC",
			MType: "gauge",
			Value: &nextGC,
		},
		{
			ID:    "NumForcedGC",
			MType: "gauge",
			Value: &numForcedGC,
		},
		{
			ID:    "NumGC",
			MType: "gauge",
			Value: &numGC,
		},
		{
			ID:    "OtherSys",
			MType: "gauge",
			Value: &otherSys,
		},
		{
			ID:    "PauseTotalNs",
			MType: "gauge",
			Value: &pauseTotalNs,
		},
		{
			ID:    "StackInuse",
			MType: "gauge",
			Value: &stackInUse,
		},
		{
			ID:    "StackSys",
			MType: "gauge",
			Value: &stackSys,
		},
		{
			ID:    "Sys",
			MType: "gauge",
			Value: &sys,
		},
		{
			ID:    "TotalAlloc",
			MType: "gauge",
			Value: &totalAlloc,
		},
		{
			ID:    "GCSys",
			MType: "gauge",
			Value: &gcSys,
		},
		{
			ID:    "HeapIdle",
			MType: "gauge",
			Value: &heapIdle,
		},
		{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &randomValue,
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: &counter,
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
