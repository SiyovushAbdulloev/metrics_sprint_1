package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	pkg_hash "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/hash"
	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/v3/mem"
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
	ConnAttempts   int
	HashKey        string
	RateLimit      int
}

func collectMetrics(m *Metrics) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	v, _ := mem.VirtualMemory()
	cpu := cpuid.CPU

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
	total := float64(v.Total)
	free := float64(v.Free)
	cpuCount := float64(cpu.LogicalCPU())

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
			ID:    "Total",
			MType: "gauge",
			Value: &total,
		},
		{
			ID:    "Free",
			MType: "gauge",
			Value: &free,
		},
		{
			ID:    "Total",
			MType: "CPU_COUNT",
			Value: &cpuCount,
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: &counter,
		},
	}

	m.data = data
}

func sendMetrics(client http.Client, ms []Metric, cfg Config) {
	for _, metric := range ms {
		var err error
		data, err := json.Marshal(metric)
		if err != nil {
			log.Printf("Error marshaling metric: %v", err)
			return
		}

		body := bytes.NewBuffer(data)
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update/", cfg.Address), body)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		if cfg.HashKey != "" {
			hash := pkg_hash.CalculateHashSHA256(body.Bytes(), cfg.HashKey)
			req.Header.Set("HashSHA256", hash)
		}

		for i := 0; i <= cfg.ConnAttempts; i++ {
			res, err2 := client.Do(req)
			err = err2
			if err2 == nil {
				res.Body.Close()
				return
			}

			if res != nil {
				res.Body.Close()
			}

			fmt.Printf("Error: %v\n", err)

			time.Sleep(time.Second * time.Duration(i*2+1))
		}

		if err != nil {
			log.Printf("Error posting metric: %v", err)
		}
	}
}

func getVars() Config {
	var address string
	var reportInterval int
	var pollInterval int
	var hashKey string
	var rateLimit int

	addr := os.Getenv("ADDRESS")
	reportInt := os.Getenv("REPORT_INTERVAL")
	pollInt := os.Getenv("POLL_INTERVAL")
	hashKeyStr := os.Getenv("KEY")
	addrFlag := flag.String("a", "localhost:8080", "The address to send HTTP requests.")
	reportIntFlag := flag.Int("r", 10, "The interval in seconds between metric reporting. (in seconds)")
	pollIntFlag := flag.Int("p", 2, "The interval in seconds between metric polling. (in seconds)")
	hashKeyFlag := flag.String("k", "", "The hash key")
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimitFlag := flag.Int("l", 5, "Max concurrent outgoing requests")
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

	if hashKeyStr == "" {
		hashKey = *hashKeyFlag
	} else {
		hashKey = hashKeyStr
	}

	if rateLimitStr == "" {
		rateLimit = *rateLimitFlag
	} else {
		val, err := strconv.Atoi(rateLimitStr)
		if err != nil {
			panic(err)
		}
		rateLimit = val
	}

	return Config{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		ConnAttempts:   3,
		HashKey:        hashKey,
		RateLimit:      rateLimit,
	}
}

type Job struct {
	Metric Metric
}

func worker(id int, jobs <-chan Job, client http.Client, cfg Config) {
	for job := range jobs {
		sendMetrics(client, []Metric{job.Metric}, cfg)
	}
}

func main() {
	config := getVars()
	fmt.Println(config)
	client := http.Client{}
	m := Metrics{}
	collectTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer collectTicker.Stop()
	defer sendTicker.Stop()

	jobs := make(chan Job, 100)

	for i := 0; i < config.RateLimit; i++ {
		go worker(i, jobs, client, config)
	}

	go func() {
		for {
			select {
			case <-collectTicker.C:
				collectMetrics(&m)
			case <-sendTicker.C:
				//go func(metrics []Metric) {
				//	sendMetrics(client, metrics, config)
				//}(append([]Metric(nil), m.data...))
				for _, metric := range m.data {
					jobs <- Job{Metric: metric}
				}
			}
		}
	}()

	select {}
}
