package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
)

// ClickhouseInspection .
func ClickhouseInspection(cfg *config.Settings, sc *schema.CheckSchema, reportDetails map[string]string) {
	limits := sc.ResourceLimits
	allAttrs := schema.MergeReportsAttrs(sc.Reports)

	log := logger.Get("check", "Clickhouse Inspection")

	// pterm.DefaultSection.Println("Clickhouse Inspection Report")

	commonInspection(log, limits, allAttrs, reportDetails)
}
