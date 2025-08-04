package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/configparam"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/crypto"
	pkg_hash "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/hash"
	pb "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/proto"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/utils/localip"
	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	pubKey     *rsa.PublicKey
	pubKeyOnce sync.Once
	pubKeyErr  error
	wg         sync.WaitGroup
)

type Build struct {
	Version string
	Date    string
	Commit  string
}

var buildInfo Build

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
	CryptoKeyPath  string
	UseGRPC        bool
}

type JSONAgentConfig struct {
	Address        string `json:"address"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

func readJSONAgentConfig(path string) (*JSONAgentConfig, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg JSONAgentConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
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
	pubKeyOnce.Do(func() {
		if cfg.CryptoKeyPath != "" {
			pubKey, pubKeyErr = crypto.LoadPublicKey(cfg.CryptoKeyPath)
		}
	})

	if pubKeyErr != nil {
		log.Printf("Public key error: %v", pubKeyErr)
		return
	}

	for _, metric := range ms {
		var err error
		data, err := json.Marshal(metric)
		if err != nil {
			log.Printf("Error marshaling metric: %v", err)
			return
		}

		if pubKey != nil {
			encrypted, err := crypto.EncryptWithPublicKey(data, pubKey)
			if err != nil {
				log.Printf("Error encrypting: %v", err)
				continue
			}
			data = encrypted
		}

		body := bytes.NewBuffer(data)
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update/", cfg.Address), body)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		realIP := localip.LocalIP()
		req.Header.Set("X-Real-IP", realIP)

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
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to JSON config file")
	flag.StringVar(&configPath, "c", "", "Path to JSON config file (short)")

	configPath = configparam.ExtractConfig()

	jsonCfg, _ := readJSONAgentConfig(configPath)
	if jsonCfg == nil {
		jsonCfg = &JSONAgentConfig{}
	}

	var address string
	var reportInterval int
	var pollInterval int
	var hashKey string
	var rateLimit int
	var cryptoKeyPath string
	var useGrpc bool

	addr := os.Getenv("ADDRESS")
	reportInt := os.Getenv("REPORT_INTERVAL")
	pollInt := os.Getenv("POLL_INTERVAL")
	hashKeyStr := os.Getenv("KEY")
	useGrpcEnv := os.Getenv("USE_GRPC")
	addrFlag := flag.String("a", "localhost:8080", "The address to send HTTP requests.")
	reportIntFlag := flag.Int("r", 10, "The interval in seconds between metric reporting. (in seconds)")
	pollIntFlag := flag.Int("p", 2, "The interval in seconds between metric polling. (in seconds)")
	hashKeyFlag := flag.String("k", "", "The hash key")
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimitFlag := flag.Int("l", 5, "Max concurrent outgoing requests")
	cryptoKeyEnv := os.Getenv("CRYPTO_KEY")
	cryptoKeyFlag := flag.String("crypto-key", "./public.pem", "Path to public key (PEM) for RSA encryption")
	useGrpcFlag := flag.Bool("grpc", false, "Use gRPC to send metrics")
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

	if cryptoKeyEnv != "" {
		cryptoKeyPath = cryptoKeyEnv
	} else {
		cryptoKeyPath = *cryptoKeyFlag
	}

	if useGrpcEnv != "" {
		useGrpc = useGrpcEnv == "true"
	} else {
		useGrpc = *useGrpcFlag
	}

	return Config{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		ConnAttempts:   3,
		HashKey:        hashKey,
		RateLimit:      rateLimit,
		CryptoKeyPath:  cryptoKeyPath,
		UseGRPC:        useGrpc,
	}
}

type Job struct {
	Metric Metric
}

func worker(id int, jobs <-chan Job, client http.Client, cfg Config, wg *sync.WaitGroup) {
	for job := range jobs {
		if cfg.UseGRPC {
			sendMetricsGRPC([]Metric{job.Metric}, cfg)
		} else {
			sendMetrics(client, []Metric{job.Metric}, cfg)
		}
		wg.Done()
	}
}

func init() {
	buildInfo.Version = "1.0.0"
	buildInfo.Date = "2025-07-04"
	buildInfo.Commit = "HEAD"
}

func sendMetricsGRPC(metrics []Metric, cfg Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("gRPC dial error: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewMetricsServiceClient(conn)

	// Преобразуем в protobuf
	var pbMetrics []*pb.Metric
	for _, m := range metrics {
		pbMetrics = append(pbMetrics, &pb.Metric{
			Id:    m.ID,
			Type:  m.MType,
			Delta: *m.Delta,
			Value: *m.Value,
		})
	}

	_, err = client.SendMetrics(ctx, &pb.MetricsRequest{
		Metrics: pbMetrics,
	})
	if err != nil {
		log.Printf("SendMetrics RPC error: %v", err)
		return
	}

	log.Println("Metrics sent via gRPC")
}

func main() {
	fmt.Printf("Build version: %s (или \"N/A\" при отсутствии значения) \n", buildInfo.Version)
	fmt.Printf("Build date: %s (или \"N/A\" при отсутствии значения) \n", buildInfo.Date)
	fmt.Printf("Build commit: %s (или \"N/A\" при отсутствии значения) \n", buildInfo.Commit)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	config := getVars()
	if config.CryptoKeyPath != "" {
		pk, err := crypto.LoadPublicKey(config.CryptoKeyPath)
		if err != nil {
			log.Printf("Failed to load public key: %v", err)
		} else {
			pubKey = pk
		}
	}
	client := http.Client{}
	m := Metrics{}
	collectTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer collectTicker.Stop()
	defer sendTicker.Stop()

	jobs := make(chan Job, 100)

	for i := 0; i < config.RateLimit; i++ {
		go worker(i, jobs, client, config, &wg)
	}

	go func() {
		<-stopChan
		log.Println("Получен сигнал завершения, завершаем агент...")

		collectTicker.Stop()
		sendTicker.Stop()
		close(jobs) // это завершит всех воркеров

		wg.Wait() // дождёмся, пока воркеры отправят всё
		os.Exit(0)
	}()

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
					wg.Add(1)
					jobs <- Job{Metric: metric}
				}
			}
		}
	}()

	select {}
}
