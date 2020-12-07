package schema

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
)

const (
	// DriverDCAP .
	DriverDCAP = "driver_dcap"
	// DriverISGX .
	DriverISGX = "driver_isgx"
	// CasTest .
	CasTest = "net_cas_test"
	// CasProd .
	CasProd = "net_cas_prod"
	// IntelPbs .
	IntelPbs = "http_intel_pbs"
	// IntelCrl .
	IntelCrl = "http_intel_pse_crl"
	// IntelOcsp .
	IntelOcsp = "http_intel_pse_ocsp"
	// IntelWl .
	IntelWl = "http_intel_wl"
)

// GetEnclaveCheckSchema .
func GetEnclaveCheckSchema(cfg *config.Settings) *CheckSchema {
	schemaReports := MergeReports(
		GetCommonSchemaReports(cfg),
		GetCommonBlockchainSchemaReports(cfg),
		GetClickhouseServicesSchemaReports(cfg),
		GetDmpEnclaveServicesSchemaReports(cfg),
		ReportsGroup{
			DriverDCAP: &reports.DriverReport{
				DriverName: "intel_sgx",
			},
			DriverISGX: &reports.DriverReport{
				DriverName: "isgx",
			},
			CasTest: &reports.NetProbeReport{
				Type:   "tcp",
				Target: cfg.Common.CasTestTarget,
			},
			CasProd: &reports.NetProbeReport{
				Type:   "tcp",
				Target: cfg.Common.CasProdTarget,
			},
			IntelPbs: &reports.HTTPReport{
				WithProxy: ".env",
				URL:       "http://ps.sgx.trustedservices.intel.com/",
			},
			IntelCrl: &reports.HTTPReport{
				WithProxy: ".env",
				URL:       "https://trustedservices.intel.com/content/CRL/",
			},
			IntelOcsp: &reports.HTTPReport{
				WithProxy: ".env",
				URL:       "http://trustedservices.intel.com/ocsp",
			},
			IntelWl: &reports.HTTPReport{
				WithProxy: ".env",
				URL:       "http://whitelist.trustedservices.intel.com/SGX/LCWL/Linux/sgx_white_list_cert.bin",
			},
		},
	)

	return &CheckSchema{
		Role: config.RoleEnclave,
		ResourceLimits: ResourceLimitsType{
			Disks + "." + reports.DisksTotalSpaceIntAttr: {
				Minimal: 200000,
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
