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
			"disk." + reports.DiskTotalSpaceIntAttr: {
				Minimal: 100000,
			},
			"cpu." + reports.CPUCoresIntAttr: {
				Minimal: 4,
			},
			"ram." + reports.RAMTotalIntAttr: {
				Minimal: 8000,
			},
		},
		Reports: schemaReports,
	}
}
