package client

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	logTimestampFormat = "2006-01-02 15:04:05"
)

// NewFileLogger creates a log file and init logger
func NewFileLogger(path string) *logrus.Logger {
	logger := logrus.New()

	// disable color
	if len(path) != 0 {
		_ = os.Setenv("NO_COLOR", "true")
	}

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: logTimestampFormat,
	}

	if file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	return logger
}
