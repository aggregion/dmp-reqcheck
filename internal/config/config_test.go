package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewSettings(t *testing.T) {
	v := viper.GetViper()
	v.Set("logging.console.level", "warning")
	v.Set("host.roles", "ch")
	settings := NewSettings(v, true)

	require.NotNil(t, settings.LoggingConsole)
	require.Equal(t, settings.LoggingConsole.Level, "warning")
}

func TestInvalidSettings(t *testing.T) {
	v := viper.GetViper()
	v.Set("logging.console.level", "invalidlevel")
	v.Set("host.roles", "ch")

	require.Panics(t, func() {
		NewSettings(v, true)
	})
}
