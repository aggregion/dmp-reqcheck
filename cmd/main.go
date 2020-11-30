package main

import (
	"fmt"

	"github.com/aggregion/dmp-reqcheck/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainCommand = &cobra.Command{
	Use:     "dmp-reqcheck",
	Long:    "dmp-reqcheck",
	Version: fmt.Sprintf("%s (git %s:%s)", internal.AppVersion, internal.GitBranch, internal.GitCommit)}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	println(internal.AppVersion)
}
