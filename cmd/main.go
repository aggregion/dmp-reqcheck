package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aggregion/dmp-reqcheck/internal"
	"github.com/aggregion/dmp-reqcheck/internal/cmds"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainCommand = &cobra.Command{
	Use:     "dmp-reqcheck",
	Long:    "dmp-reqcheck",
	Version: fmt.Sprintf("%s (git %s:%s)", internal.AppVersion, internal.GitBranch, internal.GitCommit)}

func main() {
	viper.SetEnvPrefix("drc")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	mainFlags := mainCommand.PersistentFlags()
	// mainFlags.String("config-file", "", "Config file")

	mainFlags.String("logging-level", "info", "Default console level")
	viper.BindPFlag("logging.console.level", mainFlags.Lookup("logging-level"))

	mainCommand.AddCommand(cmds.CheckCommand())
	mainCommand.AddCommand(cmds.ListenCommand())

	if mainCommand.Execute() != nil {
		os.Exit(1)
	}
}
