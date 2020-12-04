package cmds

import (
	"context"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/inspection"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CheckCommand .
func CheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check",
		Aliases: []string{"c"},
		Short:   "Do checks",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.NewSettings(viper.GetViper(), false)
			logger.InitLogger(cfg)

			inspection.RunInspection(context.Background(), cfg)
		},
	}

	flags := cmd.PersistentFlags()

	flags.String("concurrency", "", "Concurrency of reports gathering (Min 1, Max 16).")
	viper.BindPFlag("defaults.concurrency", flags.Lookup("concurrency"))

	return cmd
}
