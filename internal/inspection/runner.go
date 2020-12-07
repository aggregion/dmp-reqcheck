package inspection

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/aggregion/dmp-reqcheck/pkg/common"
	"github.com/aggregion/dmp-reqcheck/pkg/utils"
	"github.com/pterm/pterm"
)

// GetResultSchema .
func GetResultSchema(cfg *config.Settings) *schema.CheckSchema {
	allSchemas := []*schema.CheckSchema{}

	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleCH}) {
		allSchemas = append(allSchemas, schema.GetClickhouseCheckSchema(cfg))
	}
	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleDmp}) {
		allSchemas = append(allSchemas, schema.GetDmpCheckSchema(cfg))
	}
	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleEnclave}) {
		allSchemas = append(allSchemas, schema.GetEnclaveCheckSchema(cfg))
	}

	wholeSchema := schema.MergeSchemas(allSchemas...)

	return &wholeSchema
}

// RunInspection .
func RunInspection(ctx context.Context, cfg *config.Settings) error {
	wholeSchema := GetResultSchema(cfg)

	log := logger.Get("inspection", "RunInpsection")

	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("DMP", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("ReqCheck", pterm.NewStyle(pterm.FgLightRed))).
		Render()
	pterm.DefaultHeader.WithMargin(20).Printf(
		"Selected roles: %s", cfg.Host.Roles)
	pterm.Println("Hardware Resource limits are summarized")

	spinnerPrefix := fmt.Sprintf("Gathering from %d reports...", len(wholeSchema.Reports))
	spinner, _ := pterm.DefaultSpinner.WithDelay(time.Millisecond * 100).Start(spinnerPrefix)

	waitGroup := sync.WaitGroup{}
	workLimits := common.NewSemaphore(cfg.Host.GatherConcurrency)

	waitGroup.Add(len(wholeSchema.Reports))

	for _, report := range wholeSchema.Reports {
		go func(report reports.IReport) {
			workLimits.EnterWithCtx(ctx, 1)
			defer workLimits.LeaveWithCtx(ctx, 1)
			defer waitGroup.Done()

			typeStr := reflect.TypeOf(report).String()

			log.Tracef("Gather %s report...", typeStr)

			spinner.UpdateText(fmt.Sprintf("%s %s", spinnerPrefix, typeStr))

			report.Gather(ctx)

			time.Sleep(time.Millisecond * 100)

			log.Debugf("Gathered %s:%+v report.", typeStr, report)
		}(report)
	}

	waitGroup.Wait()

	spinner.Success("Reports are gathered.")
	spinner.Stop()

	reportsDetails := map[string]string{}
	for name, report := range wholeSchema.Reports {
		reportsDetails[name] = report.String()
	}

	log.Debug("Match reports with specified Roles...")

	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleCH}) {
		ClickhouseInspection(cfg, wholeSchema, reportsDetails)
	}
	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleDmp}) {
		DmpInspection(cfg, wholeSchema, reportsDetails)
	}
	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleEnclave}) {
		EnclaveInspection(cfg, wholeSchema, reportsDetails)
	}

	pterm.Println()

	log.Debug("All Checks are Complete!")

	return nil
}
