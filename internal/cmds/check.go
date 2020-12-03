package cmds

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/inspection"
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
			cfg := config.NewSettings(viper.GetViper())

			inspection.RunInspection(cfg)
		},
	}

	flags := cmd.PersistentFlags()

	flags.String("roles", "", "Roles")
	viper.BindPFlag("host.roles", flags.Lookup("roles"))

	flags.String("hosts", "", "Hosts")
	viper.BindPFlag("host.hosts", flags.Lookup("hosts"))

	return cmd
}
