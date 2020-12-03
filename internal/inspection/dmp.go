package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
)

// DmpInspection .
func DmpInspection(cfg *config.Settings, sc *schema.CheckSchema) {
	limits := sc.ResourceLimits
	allAttrs := schema.MergeReportsAttrs(sc.Reports)

	log := logger.Get("check", "DMP Inspection")

	commonInspection(log, limits, allAttrs)
}
