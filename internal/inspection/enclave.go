package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
)

// EnclaveInspection .
func EnclaveInspection(cfg *config.Settings, sc *schema.CheckSchema) {
	limits := sc.ResourceLimits
	allAttrs := schema.MergeReportsAttrs(sc.Reports)

	log := logger.Get("check", "Enclave Inspection")

	commonInspection(log, limits, allAttrs)

	// var intVal, intVal2, minVal int64
	var intVal int64
	var strVal string

	//
	// Hardware
	//
	intVal = reportIntAttr(allAttrs, "cpu", reports.CPUSgx1IntAttr)
	if intVal != 1 {
		log.Error("The CPU hasn't SGX1 or it not enabled")
	}

	intVal = reportIntAttr(allAttrs, "cpu", reports.CPUSgxFclIntAttr)
	if intVal != 1 {
		log.Warn("The CPU hasn't FCL feature or it not enabled, not possible to install Intel DCAP driver")
	}

	strVal = reportStrAttr(allAttrs, "os", reports.OSVendorStrAttr)
	if strVal != "ubuntu" && strVal != "centos" {
		log.Warn("The OS is %s, the driver might be not install properly, should be use Ubuntu 18.x or CentOS 8.x", strVal)
	}

	strVal = reportStrAttr(allAttrs, "os", reports.OSVersionStrAttr)
	intVal = reportIntAttr(allAttrs, "os", reports.OSMajorVersionIntAttr)
	switch strVal {
	case "ubuntu":
		if intVal < 18 {
			log.Warningf("The Ubuntu major version %s is too old, minimal version is 18.x", strVal)
		}
	case "centos":
		if intVal < 8 {
			log.Warningf("The CentOS major version %s is too old, minimal version is 8.x", strVal)
		}
	}

	// intVal = reportIntAttr(allAttrs, "cpu", reports.HTTPBodyStrAttr)
	// if intVal != 1 {
	// 	log.Warn("The CPU hasn't FCL feature or it not enabled, not possible to install Intel DCAP driver")
	// }
}
