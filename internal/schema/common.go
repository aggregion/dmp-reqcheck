package schema

import (
	"fmt"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/pkg/utils"
)

// GetCommonSchemaReports .
func GetCommonSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		"cpu":        &reports.CPUReport{},
		"ram":        &reports.RAMReport{},
		"disk":       &reports.DiskReport{},
		"os":         &reports.OSReport{},
		"kernel":     &reports.KernelReport{},
		"hypervisor": &reports.HypervisorReport{},

		"docker": &reports.DockerReport{},
		"docker_registry": &reports.HTTPReport{
			URL: "https://registry.aggregion.com",
		},
		"net_aggregion_proxy": &reports.NetProbeReport{
			Type:   "tcp",
			Target: "185.175.44.42:80",
		},
	}
}

// GetCommonBlockchainSchemaReports .
func GetCommonBlockchainSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		"net_blockchain_testnet": &reports.NetProbeReport{
			Type:   "tcp",
			Target: "185.137.232.118:9999",
		},
		"net_blockchain_prodnet": &reports.NetProbeReport{
			Type:   "tcp",
			Target: "185.137.232.118:8888",
		},
	}
}

// GetClickhouseServicesSchemaReports .
func GetClickhouseServicesSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		"host_service_clickhouse": &reports.HostServiceReport{
			IsThisHost: utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleCH}),
			Target:     fmt.Sprintf("http://%s:%d", cfg.Host.Hosts[config.RoleCH], cfg.Host.DefaultClickhousePort),
		},
	}
}

// GetDmpEnclaveServicesSchemaReports .
func GetDmpEnclaveServicesSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		"host_service_dmp": &reports.HostServiceReport{
			IsThisHost: utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleDmp}),
			Target:     fmt.Sprintf("http://%s:8080", cfg.Host.Hosts[config.RoleDmp]),
		},
		"host_service_enclave": &reports.HostServiceReport{
			IsThisHost: utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleEnclave}),
			Target:     fmt.Sprintf("http://%s:8321", cfg.Host.Hosts[config.RoleEnclave]),
		},
	}
}

// MergeReports .
func MergeReports(groups ...ReportsGroup) (out ReportsGroup) {
	out = make(map[string]reports.IReport)

	for _, group := range groups {
		for name, report := range group {
			out[name] = report
		}
	}

	return
}

// MergeReportsAttrs .
func MergeReportsAttrs(reportItems ReportsGroup) (out map[string]interface{}) {
	out = make(map[string]interface{})

	for name, report := range reportItems {
		attrs := reports.GetReportAttrs(report)
		for attrName, attrValue := range attrs {
			out[name+"."+attrName] = attrValue
		}
	}

	return
}

// MergeResourceLimits .
func MergeResourceLimits(groups ...ResourceLimitsType) (out ResourceLimitsType) {
	out = make(ResourceLimitsType)

	for _, group := range groups {
		for name, limits := range group {
			var newLimits ResourceLimitType

			if current, ok := out[name]; ok {
				newLimits = ResourceLimitType{
					Minimal:   limits.Minimal + current.Minimal,
					Recommend: limits.Recommend + current.Recommend,
				}
				out[name] = newLimits
			} else {
				out[name] = limits
			}
		}
	}

	return
}

// MergeSchemas .
func MergeSchemas(schemas ...*CheckSchema) (out CheckSchema) {
	limits := make(ResourceLimitsType)
	reps := make(ReportsGroup)

	for _, schema := range schemas {
		limits = MergeResourceLimits(limits, schema.ResourceLimits)
		reps = MergeReports(reps, schema.Reports)
	}

	out.ResourceLimits = limits
	out.Reports = reps

	return
}
