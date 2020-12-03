package schema

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
)

// GetDmpCheckSchema .
func GetDmpCheckSchema(cfg *config.Settings) *CheckSchema {
	schemaReports := MergeReports(
		GetCommonSchemaReports(cfg),
		GetCommonBlockchainSchemaReports(cfg),
		GetClickhouseServicesSchemaReports(cfg),
		GetDmpEnclaveServicesSchemaReports(cfg),
	)

	return &CheckSchema{
		Role: config.RoleDmp,
		ResourceLimits: ResourceLimitsType{
			"disk." + reports.DiskTotalSpaceIntAttr: {
				Minimal: 200000,
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
