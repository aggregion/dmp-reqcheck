package cmds

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/inspection"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListenCommand .
func ListenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "listen",
		Aliases: []string{"l"},
		Short:   "Listen mode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.NewSettings(viper.GetViper(), true)
			logger.InitLogger(cfg)

			ctx, stop := context.WithCancel(context.Background())
			defer stop()

			schema := inspection.GetResultSchema(cfg)

			for _, report := range schema.Reports {
				report.Start(ctx)
			}

			appSignals := make(chan os.Signal, 1)
			signal.Notify(appSignals,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT,
			)
			go func() {
				sig := <-appSignals
				fmt.Println("Stop process. Catch signal", sig)
				stop()
			}()

			fmt.Println("DMP-ReqCheck. Selected Roles", cfg.Host.Roles)
			fmt.Println("Listen... Press CTRL-C to stop")
			<-ctx.Done()

			for _, report := range schema.Reports {
				report.Stop(ctx)
			}
		},
	}

	return cmd
}
