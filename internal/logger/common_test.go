package logger

import (
	"testing"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestInitCommon(t *testing.T) {
	initCommon(&config.Settings{
		LoggingConsole: &config.LoggingConsoleSettings{
			Level: "error",
		},
	})

	require.Equal(t, GetRootLogger().Level, logrus.ErrorLevel, "expected logging level to be error")
}
