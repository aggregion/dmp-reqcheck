package schema

import (
	"fmt"

	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/pkg/utils"
)

const (
	// CPU .
	CPU = "cpu"
	// RAM .
	RAM = "ram"
	// Disks .
	Disks = "disks"

	// OS .
	OS = "os"
	// Kernel .
	Kernel = "kernel"
	// HV .
	HV = "hypervisor"

	// Podman .
	Podman = "podman"
	// Docker .
	Docker = "docker"
	// DockerRegistry .
	DockerRegistry = "docker_registry"
	// AggregionProxy .
	AggregionProxy = "net_aggregion_proxy"
	// EOSTestNet .
	EOSTestNet = "net_blockchain_testnet"
	// EOSProdNet .
	EOSProdNet = "net_blockchain_prodnet"
	// DmpHostSvc .
	DmpHostSvc = "host_service_dmp"
	// EnclaveHostSvc .
	EnclaveHostSvc = "host_service_enclave"
	// ClickhouseHostSvc .
	ClickhouseHostSvc = "host_service_clickhouse"
)

// GetCommonSchemaReports .
func GetCommonSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		CPU:    &reports.CPUReport{},
		RAM:    &reports.RAMReport{},
		Disks:  &reports.DisksReport{},
		OS:     &reports.OSReport{},
		Kernel: &reports.KernelReport{},
		HV:     &reports.HypervisorReport{},

		Docker: &reports.DockerReport{},
		Podman: &reports.PodmanReport{},
		DockerRegistry: &reports.HTTPReport{
			WithProxy: ".env",
			URL:       "https://registry.aggregion.com",
		},
		AggregionProxy: &reports.NetProbeReport{
			Type:   "tcp",
			Target: "185.175.44.42:80",
		},
	}
}

// GetCommonBlockchainSchemaReports .
func GetCommonBlockchainSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		EOSTestNet: &reports.NetProbeReport{
			Type:   "tcp",
			Target: "185.137.232.118:9999",
		},
		EOSProdNet: &reports.NetProbeReport{
			Type:   "tcp",
			Target: "185.137.232.118:8888",
		},
	}
}

// GetClickhouseServicesSchemaReports .
func GetClickhouseServicesSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		ClickhouseHostSvc: &reports.HostServiceReport{
			IsThisHost: utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleCH}),
			Target:     fmt.Sprintf("http://%s:%d", cfg.Host.Hosts[config.RoleCH], cfg.Host.DefaultClickhousePort),
		},
	}
}

// GetDmpEnclaveServicesSchemaReports .
func GetDmpEnclaveServicesSchemaReports(cfg *config.Settings) ReportsGroup {
	return ReportsGroup{
		DmpHostSvc: &reports.HostServiceReport{
			IsThisHost: utils.IsIntersectStrs(cfg.Host.Roles, config.Roles{config.RoleDmp}),
			Target:     fmt.Sprintf("http://%s:8080", cfg.Host.Hosts[config.RoleDmp]),
		},
		EnclaveHostSvc: &reports.HostServiceReport{
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
