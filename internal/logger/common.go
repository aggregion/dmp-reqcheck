package logger

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/sirupsen/logrus"
)

var rootLogger = logrus.New()

// Get .
func Get(tag, value string) *logrus.Entry {
	return rootLogger.WithFields(logrus.Fields{tag: value})
}

// GetRootLogger .
func GetRootLogger() *logrus.Logger {
	return rootLogger
}

func initCommon(settings *config.Settings) {
	switch settings.LoggingConsole.Level {
	case "trace":
		rootLogger.SetLevel(logrus.TraceLevel)
	case "debug":
		rootLogger.SetLevel(logrus.DebugLevel)
	case "warning":
		rootLogger.SetLevel(logrus.WarnLevel)
	case "info":
		rootLogger.SetLevel(logrus.InfoLevel)
	case "error":
		rootLogger.SetLevel(logrus.ErrorLevel)
	default:
		rootLogger.SetLevel(logrus.InfoLevel)
	}
}

// InitLogger .
func InitLogger(settings *config.Settings) {
	initCommon(settings)
}
