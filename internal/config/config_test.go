package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewSettings(t *testing.T) {
	v := viper.GetViper()
	v.Set("logging.console.level", "warning")
	settings := NewSettings(v)

	require.NotNil(t, settings.LoggingConsole)
	require.Equal(t, settings.LoggingConsole.Level, "warning")
}

func TestInvalidSettings(t *testing.T) {
	v := viper.GetViper()
	v.Set("logging.console.level", "invalidlevel")

	require.Panics(t, func() {
		NewSettings(v)
	})
}
