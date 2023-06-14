package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

const logTimestampFormat = "2006-01-02 15:04:05"

func Level(env string) logrus.Level {
	switch env {
	case "test":
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}

// NewFileLogger creates a log file and init logger
func NewFileLogger(path string, env string) *logrus.Logger {
	logger := logrus.New()

	// default log to os.Stderr
	if path == "" {
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: logTimestampFormat,
		}
		return logger
	}

	if file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
	logger.Level = logrus.DebugLevel
	logger.ReportCaller = true
	// use json formatter
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: logTimestampFormat,
	}
	//	logger.WithTime(time.Now().Local())
	return logger
}
