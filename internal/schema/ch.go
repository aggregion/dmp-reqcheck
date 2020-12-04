package schema

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
)

// GetClickhouseCheckSchema .
func GetClickhouseCheckSchema(cfg *config.Settings) *CheckSchema {
	schemaReports := MergeReports(
		GetCommonSchemaReports(cfg),
		GetClickhouseServicesSchemaReports(cfg),
	)

	return &CheckSchema{
		Role: config.RoleCH,
		ResourceLimits: ResourceLimitsType{
			Disks + "." + reports.DisksTotalSpaceIntAttr: {
				Minimal: 100000,
			},
			CPU + "." + reports.CPUCoresIntAttr: {
				Minimal: 4,
			},
			RAM + "." + reports.RAMTotalIntAttr: {
				Minimal: 8000,
			},
		},
		Reports: schemaReports,
	}
}
