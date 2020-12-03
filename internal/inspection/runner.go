package inspection

import (
	"context"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/aggregion/dmp-reqcheck/pkg/utils"
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

	for _, report := range wholeSchema.Reports {
		report.Gather(context.Background())
	}

	return &wholeSchema
}

// RunInspection .
func RunInspection(cfg *config.Settings) {
	wholeSchema := GetResultSchema(cfg)

	log := logger.Get("inspection", "RunInpsection")

	log.Infof("Selected roles: %s", cfg.Host.Roles)

	log.Infof("Gathering from %d reports...", len(wholeSchema.Reports))
	for _, report := range wholeSchema.Reports {
		report.Gather(context.Background())
	}

	log.Info("Match reports with specified Roles...")

	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleCH}) {
		ClickhouseInspection(cfg, wholeSchema)
	}
	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleDmp}) {
		DmpInspection(cfg, wholeSchema)
	}
	if utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleEnclave}) {
		EnclaveInspection(cfg, wholeSchema)
	}

	log.Info("All Checks are Complete!")
}
