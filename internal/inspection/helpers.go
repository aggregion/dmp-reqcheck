package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/sirupsen/logrus"
)

func reportIntAttr(attrs map[string]interface{}, repID string, attrName string) int64 {
	if val, ok := attrs[repID+"."+attrName]; ok {
		return val.(int64)
	}

	return 0
}

func reportStrAttr(attrs map[string]interface{}, repID string, attrName string) string {
	if val, ok := attrs[repID+"."+attrName]; ok {
		return val.(string)
	}

	return ""
}

func getLimit(attrs schema.ResourceLimitsType, repID string, attrName string) schema.ResourceLimitType {
	if val, ok := attrs[repID+"."+attrName]; ok {
		return val
	}

	return schema.ResourceLimitType{}
}

// commonInspection .
func commonInspection(log *logrus.Entry, limits schema.ResourceLimitsType, allAttrs map[string]interface{}) {
	var intVal, minVal int64
	var strVal string

	//
	// Hardware
	//
	intVal = reportIntAttr(allAttrs, "cpu", reports.CPUCoresIntAttr)
	minVal = getLimit(limits, "cpu", reports.CPUCoresIntAttr).Minimal

	if intVal < minVal {
		log.Warningf("The CPU has only %d cores, minimum count is %d", intVal, minVal)
	}

	intVal = reportIntAttr(allAttrs, "disk", reports.DiskTotalSpaceIntAttr)
	minVal = getLimit(limits, "disk", reports.DiskTotalSpaceIntAttr).Minimal
	if intVal < minVal {
		log.Warningf("The Disk has only %d MB, minimal size is %d MB", intVal, minVal)
	}

	intVal = reportIntAttr(allAttrs, "ram", reports.RAMTotalIntAttr)
	minVal = getLimit(limits, "ram", reports.RAMTotalIntAttr).Minimal
	if intVal < minVal {
		log.Warningf("The RAM has only %d MB, minimal size is %d MB", intVal, minVal)
	}

	//
	// OS
	//
	strVal = reportStrAttr(allAttrs, "os", reports.OSVendorStrAttr)
	if strVal != "ubuntu" && strVal != "centos" {
		log.Infof("The OS is %s, it would be better to use Ubuntu 18.x or CentOS 8.x", strVal)
	}

	strVal = reportStrAttr(allAttrs, "os", reports.OSVersionStrAttr)
	intVal = reportIntAttr(allAttrs, "os", reports.OSMajorVersionIntAttr)
	switch strVal {
	case "ubuntu":
		if intVal < 18 {
			log.Warningf("The Ubuntu major version %s is too old, minimal version is 18.x", strVal)
		}
	case "centos":
		if intVal < 18 {
			log.Warningf("The CentOS major version %s is too old, minimal version is 8.x", strVal)
		}
	case "rhel":
		if intVal < 18 {
			log.Warningf("The CentOS major version %s is too old, minimal version is 8.x", strVal)
		}
	}

	//
	// Resources
	//
	intVal = reportIntAttr(allAttrs, "docker", reports.DockerClientMajorVersionIntAttr)
	if intVal == 0 {
		log.Error("The Docker is not installed")
	} else if intVal < 19 {
		log.Warningf("The Docker version %d is too old, minimal version is 19.x", intVal)
	}

	intVal = reportIntAttr(allAttrs, "docker_registry", reports.HTTPStatusIntAttr)
	if intVal != 200 {
		log.Warn("The Aggregion Docker Registry (https://registry.aggregion.com) is not accessible")
	}
}
