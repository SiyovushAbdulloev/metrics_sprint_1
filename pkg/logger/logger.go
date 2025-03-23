package logger

import (
	"go.uber.org/zap"
)

type Interface interface {
	Info(message string, args ...interface{})
}

type Logger struct {
	logger *zap.SugaredLogger
}

func New() (*Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger: logger.Sugar(),
	}, nil
}

func (l Logger) Info(message string, args ...interface{}) {
	l.logger.Infow(message, args...)
}
