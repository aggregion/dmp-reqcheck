package config

import (
	"github.com/spf13/viper"
)

type (
	// Settings .
	Settings struct {
		LoggingConsole *LoggingConsoleSettings
	}
)

// NewSettings .
func NewSettings(v *viper.Viper) *Settings {
	return &Settings{
		// Loggins
		LoggingConsole: loggingConsoleSettingsValidateAndGet(v),
	}
}
