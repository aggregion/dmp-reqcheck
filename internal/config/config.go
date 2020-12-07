package config

import (
	"github.com/spf13/viper"
)

type (
	// Settings .
	Settings struct {
		LoggingConsole *LoggingConsoleSettings
		Host           *HostSettings
		Common         *CommonSettings
	}
)

// NewSettings .
func NewSettings(v *viper.Viper, isListenContext bool) *Settings {
	return &Settings{
		// Loggins
		LoggingConsole: loggingConsoleSettingsValidateAndGet(v),
		// Host
		Host: hostSettingsValidateAndGet(v, isListenContext),
		// Common
		Common: commonSettingsValidateAndGet(v),
	}
}
