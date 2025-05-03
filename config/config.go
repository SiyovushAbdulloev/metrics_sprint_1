package config

import (
	"flag"
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
}

type Server struct {
	Address string
}

type Log struct {
	Level string
}

func New() (*Config, error) {
	var address string
	var logLevel string
	var restore bool
	var storeInterval int
	var filePath string
	var dsn string
	var hashKey string

	addr := os.Getenv("ADDRESS")
	addrFlag := flag.String("a", "localhost:8080", "The address to listen on for HTTP requests.")

	ll := os.Getenv("LOG_LEVEL")
	logLevelFlag := flag.String("ll", "info", "The log level to use")

	rest := os.Getenv("RESTORE")
	restFlag := flag.Bool("r", false, "Restore or no saved data after server start")

	storeInt := os.Getenv("STORE_INTERVAL")
	storeIntFlag := flag.Int("i", 300, "After certain seconds current data will be stored in a file.")

	fp := os.Getenv("FILE_STORAGE_PATH")
	fpFlag := flag.String("f", "storage.txt", "The filepath where will be stored data from storage.")

	db := os.Getenv("DATABASE_DSN")
	//postgres://postgres:postgres@localhost:5432/metrics
	dbFlag := flag.String("d", "", "The dsn of postgresql.")

	hk := os.Getenv("KEY")
	//postgres://postgres:postgres@localhost:5432/metrics
	hkFlag := flag.String("k", "", "The hash key.")

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

	return &Config{
		Server: Server{
			Address: address,
		},
		Log: Log{
			Level: logLevel,
		},
		App: App{
			StoreInterval: storeInterval,
			Filepath:      filePath,
			Restore:       restore,
			HashKey:       hashKey,
		},
		Database: Database{
			DSN: dsn,
		},
	}, nil
}
