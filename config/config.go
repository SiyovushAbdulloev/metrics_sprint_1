package config

import (
	"flag"
	"os"
)

type Config struct {
	Server Server
	Log    Log
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

	addr := os.Getenv("ADDRESS")
	addrFlag := flag.String("a", "localhost:8080", "The address to listen on for HTTP requests.")

	ll := os.Getenv("LOG_LEVEL")
	logLevelFlag := flag.String("ll", "info", "The log level to use")

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

	return &Config{
		Server: Server{
			Address: address,
		},
		Log: Log{
			Level: logLevel,
		},
	}, nil
}
