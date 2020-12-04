package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/pterm/pterm"
)

// DmpInspection .
func DmpInspection(cfg *config.Settings, sc *schema.CheckSchema, reportDetails map[string]string) {
	limits := sc.ResourceLimits
	allAttrs := schema.MergeReportsAttrs(sc.Reports)

	log := logger.Get("check", "DMP Inspection")

	// pterm.DefaultSection.Println("DMP Inspection Report")

	commonInspection(log, limits, allAttrs, reportDetails)
	commonBlockchainInspection(log, limits, allAttrs, reportDetails)

	pterm.DefaultSection.Println("DMP Specific Report")

	var intVal int64
	// var strVal string

	intVal = reportIntAttr(allAttrs, schema.ClickhouseHostSvc, reports.HostAccessibleIntAttr)
	if intVal != 1 {
		pterm.Error.Printf("Clickhouse Host: the host %s is not accessible or not response\n", reportDetails[schema.ClickhouseHostSvc])
	} else {
		pterm.Success.Println("Clickhouse Host: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.EnclaveHostSvc, reports.HostAccessibleIntAttr)
	if intVal != 1 {
		pterm.Error.Printf("Enclave Host: the host %s is not accessible or not response\n", reportDetails[schema.EnclaveHostSvc])
	} else {
		pterm.Success.Println("Enclave Host: OK")
	}
}
