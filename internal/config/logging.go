package config

import (
	"github.com/aggregion/dmp-reqcheck/pkg/utils"
	"github.com/spf13/viper"
)

type (
	// LoggingConsoleSettings .
	LoggingConsoleSettings struct {
		Level string `validate:"required,oneof=debug info warning error"`
	}
)

func loggingConsoleSettingsValidateAndGet(v *viper.Viper) *LoggingConsoleSettings {
	var conf = &LoggingConsoleSettings{
		Level: v.GetString("logging.console.level"),
	}

	utils.MustValidate(conf)

	return conf
}
