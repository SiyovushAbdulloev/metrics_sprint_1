package config

import (
	"encoding/json"
	"flag"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/configparam"
	"os"
	"strconv"
)

type Config struct {
	Server   Server
	Log      Log
	App      App
	Database Database
}

type Database struct {
	DSN string
}

type App struct {
	StoreInterval int
	Filepath      string
	Restore       bool
	HashKey       string
	CryptoKeyPath string
}

type Server struct {
	Address       string
	TrustedSubnet string
}

type Log struct {
	Level string
}

type JSONConfig struct {
	Address       string `json:"address"`
	Restore       *bool  `json:"restore"`
	StoreInterval int    `json:"store_interval"`
	Filepath      string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
	HashKey       string `json:"hash_key"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`
}

func readJSONConfig(path string) (*JSONConfig, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg JSONConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func New() (*Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to JSON config file")
	flag.StringVar(&configPath, "c", "", "Path to JSON config file (short)")

	configPath = configparam.ExtractConfig()

	jsonCfg, _ := readJSONConfig(configPath)
	if jsonCfg == nil {
		jsonCfg = &JSONConfig{}
	}

	var address string
	var logLevel string
	var restore bool
	var storeInterval int
	var filePath string
	var dsn string
	var hashKey string
	var cryptoKeyPath string
	var trustedSubnet string

	addr := os.Getenv("ADDRESS")
	addrFlag := flag.String("a", getString(jsonCfg.Address, "localhost:8080"), "The address to listen on for HTTP requests.")

	ll := os.Getenv("LOG_LEVEL")
	logLevelFlag := flag.String("ll", "info", "The log level to use")

	rest := os.Getenv("RESTORE")
	restFlag := flag.Bool("r", getBool(jsonCfg.Restore, false), "Restore or no saved data after server start")

	storeInt := os.Getenv("STORE_INTERVAL")
	storeIntFlag := flag.Int("i", getInt(jsonCfg.StoreInterval, 300), "After certain seconds current data will be stored in a file.")

	fp := os.Getenv("FILE_STORAGE_PATH")
	fpFlag := flag.String("f", getString(jsonCfg.Filepath, "storage.txt"), "The filepath where will be stored data from storage.")

	db := os.Getenv("DATABASE_DSN")
	//postgres://postgres:postgres@localhost:5432/metrics
	dbFlag := flag.String("d", getString(jsonCfg.DatabaseDSN, ""), "The dsn of postgresql.")

	hk := os.Getenv("KEY")
	//postgres://postgres:postgres@localhost:5432/metrics
	hkFlag := flag.String("k", "", "The hash key.")

	ck := os.Getenv("CRYPTO_KEY")
	ckFlag := flag.String("crypto-key", getString(jsonCfg.CryptoKey, "./private.pem"), "Path to RSA private key file")

	trSubnet := os.Getenv("TRUSTED_SUBNET")
	trSubnetFlag := flag.String("t", getString(jsonCfg.TrustedSubnet, ""), "List of subnets, which are allowed to make request")

	flag.Parse()

	if addr == "" {
		address = *addrFlag
	} else {
		address = addr
	}

	if ll == "" {
		logLevel = *logLevelFlag
	} else {
		logLevel = ll
	}

	if rest == "" {
		restore = *restFlag
	} else {
		restore = rest == "true"
	}

	if storeInt == "" {
		storeInterval = *storeIntFlag
	} else {
		value, err := strconv.Atoi(storeInt)
		if err != nil {
			return nil, err
		}
		storeInterval = value
	}

	if fp == "" {
		filePath = *fpFlag
	} else {
		filePath = fp
	}

	if db == "" {
		dsn = *dbFlag
	} else {
		dsn = db
	}

	if hk == "" {
		hashKey = *hkFlag
	} else {
		hashKey = hk
	}

	if ck == "" {
		cryptoKeyPath = *ckFlag
	} else {
		cryptoKeyPath = ck
	}

	if trSubnet == "" {
		trustedSubnet = *trSubnetFlag
	} else {
		trustedSubnet = trSubnet
	}

	return &Config{
		Server: Server{
			Address:       address,
			TrustedSubnet: trustedSubnet,
		},
		Log: Log{
			Level: logLevel,
		},
		App: App{
			StoreInterval: storeInterval,
			Filepath:      filePath,
			Restore:       restore,
			HashKey:       hashKey,
			CryptoKeyPath: cryptoKeyPath,
		},
		Database: Database{
			DSN: dsn,
		},
	}, nil
}

func getString(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func getInt(value, defaultValue int) int {
	if value != 0 {
		return value
	}
	return defaultValue
}

func getBool(val *bool, fallback bool) bool {
	if val != nil {
		return *val
	}
	return fallback
}
