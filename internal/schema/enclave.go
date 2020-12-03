package schema

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
)

// GetEnclaveCheckSchema .
func GetEnclaveCheckSchema(cfg *config.Settings) *CheckSchema {
	schemaReports := MergeReports(
		GetCommonSchemaReports(cfg),
		GetCommonBlockchainSchemaReports(cfg),
		GetClickhouseServicesSchemaReports(cfg),
		GetDmpEnclaveServicesSchemaReports(cfg),
		ReportsGroup{
			"driver_dcap": &reports.DriverReport{
				DriverName: "intel_sgx",
			},
			"driver_isgx": &reports.DriverReport{
				DriverName: "isgx",
			},
			"net_cas_test": &reports.NetProbeReport{
				Type:   "tcp",
				Target: "185.175.44.42:18765",
			},
			"net_cas_prod": &reports.NetProbeReport{
				Type:   "tcp",
				Target: "185.175.44.40:18765",
			},
			"http_intel_pbs": &reports.HTTPReport{
				URL: "http://ps.sgx.trustedservices.intel.com/",
			},
			"http_intel_pse_crl": &reports.HTTPReport{
				URL: "https://trustedservices.intel.com/content/CRL/",
			},
			"http_intel_pse_ocsp": &reports.HTTPReport{
				URL: "http://trustedservices.intel.com/ocsp",
			},
			"http_intel_wl": &reports.HTTPReport{
				URL: "http://whitelist.trustedservices.intel.com/SGX/LCWL/Linux/sgx_white_list_cert.bin",
			},
		},
	)

	return &CheckSchema{
		Role: config.RoleEnclave,
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
