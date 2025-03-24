package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Server Server
	Log    Log
	App    App
}

type App struct {
	StoreInterval int
	Filepath      string
	Restore       bool
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
		},
	}, nil
}
