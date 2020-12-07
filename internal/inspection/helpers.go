package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/pterm/pterm"
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

var commonInspecionWasDone = false

// commonInspection .
func commonInspection(log *logrus.Entry, limits schema.ResourceLimitsType, allAttrs map[string]interface{}, reportDetails map[string]string) {
	if commonInspecionWasDone {
		return
	}

	commonInspecionWasDone = true

	var intVal, intVal2, minVal int64
	var strVal string

	pterm.DefaultSection.Println("Hardware")

	//
	// Hardware
	//
	intVal = reportIntAttr(allAttrs, schema.CPU, reports.CPUCoresIntAttr)
	minVal = getLimit(limits, schema.CPU, reports.CPUCoresIntAttr).Minimal

	if intVal < minVal {
		pterm.Warning.Printf("CPU: has only %d cores, minimum count is %d\n", intVal, minVal)
	} else {
		pterm.Success.Println("CPU: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.Disks, reports.DisksTotalSpaceIntAttr)
	minVal = getLimit(limits, schema.Disks, reports.DisksTotalSpaceIntAttr).Minimal
	if intVal < minVal {
		pterm.Warning.Printf("Disks: has only %d MB of total space, minimal size is %d MB\n", intVal, minVal)
	} else {
		pterm.Success.Println("Disks: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.RAM, reports.RAMTotalIntAttr)
	minVal = getLimit(limits, schema.RAM, reports.RAMTotalIntAttr).Minimal
	if intVal < minVal {
		pterm.Warning.Printf("RAM: has only %d MB of total size, minimal size is %d MB\n", intVal, minVal)
	} else {
		pterm.Success.Println("RAM: OK")
	}

	//
	// OS
	//
	pterm.DefaultSection.Println("OS")

	osVendor := reportStrAttr(allAttrs, schema.OS, reports.OSVendorStrAttr)
	if osVendor != "ubuntu" && osVendor != "centos" && osVendor != "rhel" && osVendor != "redhat" {
		pterm.Info.Printf("Vendor: current vendor is %s, it would be better to use Ubuntu 18.x or CentOS 8.x\n", osVendor)
	} else {
		pterm.Success.Println("Vendor: OK")
	}

	osVersionOk := true
	strVal = reportStrAttr(allAttrs, schema.OS, reports.OSVersionStrAttr)
	intVal = reportIntAttr(allAttrs, schema.OS, reports.OSMajorVersionIntAttr)
	switch osVendor {
	case "ubuntu":
		if intVal < 18 {
			pterm.Warning.Printf("The Ubuntu major version %s is too old, minimal version is 18.x\n", strVal)
			osVersionOk = false
		}
	case "centos":
		if intVal < 8 {
			pterm.Warning.Printf("The CentOS major version %s is too old, minimal version is 8.x\n", strVal)
			osVersionOk = false
		}
	case "redhat":
	case "rhel":
		if intVal < 8 {
			pterm.Warning.Printf("The RedHat major version %s is too old, minimal version is 8.x\n", strVal)
			osVersionOk = false
		}
	default:
		pterm.Warning.Printf("Detected OS is %s %s, for better compatibles use Ubuntu 18 or CentOS 8\n", osVendor, strVal)
	}
	if osVersionOk {
		pterm.Success.Println("Version: OK")
	}

	strVal = reportStrAttr(allAttrs, schema.Kernel, reports.KernelVersionStrAttr)
	intVal = reportIntAttr(allAttrs, schema.Kernel, reports.KernelMajorVersionIntAttr)
	if intVal == 0 {
		pterm.Error.Println("Kernel Version: fail to determinate version")
	} else if intVal < 4 {
		pterm.Warning.Printf("Kernel Version: major version %s is too old, minimal version is 4.x\n", strVal)
	} else {
		pterm.Success.Println("Kernel Version: OK")
	}

	strVal = reportStrAttr(allAttrs, schema.HV, reports.HypervisorNameStrAttr)
	if strVal != "" {
		pterm.Info.Printf("Hypervisor: the Host use %s as hypervisor, it would be better to use baremetal machine\n", strVal)
	} else {
		pterm.Success.Println("No-Hypervisor: OK or Not Detected")
	}

	//
	// Resources
	//
	pterm.DefaultSection.Println("Common Resources")

	isDockerExists := false
	intVal = reportIntAttr(allAttrs, schema.Docker, reports.DockerClientMajorVersionIntAttr)
	if intVal == 0 {
		pterm.Error.Println("The Docker Client is not installed or not found in default paths")
	} else if intVal < 19 {
		pterm.Warning.Printf("The Docker Client version %d is too old, minimal version is 19.x\n", intVal)
	} else {
		pterm.Success.Println("Docker Client: OK")
		isDockerExists = true
	}

	intVal = reportIntAttr(allAttrs, schema.Docker, reports.DockerServerMajorVersionIntAttr)
	if intVal == 0 {
		pterm.Error.Println("The Docker Daemon is not installed or not found in default paths")
	} else if intVal < 19 {
		pterm.Warning.Printf("The Docker Daemon version %d is too old, minimal version is 19.x\n", intVal)
	} else {
		pterm.Success.Println("Docker Daemon: OK")
		isDockerExists = true
	}

	intVal = reportIntAttr(allAttrs, schema.Docker, reports.DockerComposeMajorVersionIntAttr)
	intVal2 = reportIntAttr(allAttrs, schema.Docker, reports.DockerComposeMinorVersionIntAttr)
	if intVal == 0 {
		pterm.Warning.Println("The Docker Compose is not installed or not found in default paths")
	} else if intVal < 2 && intVal2 < 20 {
		pterm.Warning.Printf("The Docker Compose version %d.%d is too old, minimal version is 1.20.x\n", intVal, intVal2)
	} else {
		pterm.Success.Println("Docker Compose: OK")
	}

	if !isDockerExists {
		intVal = reportIntAttr(allAttrs, schema.Podman, reports.PodmanMajorVersionIntAttr)
		intVal2 = reportIntAttr(allAttrs, schema.Podman, reports.PodmanMinorVersionIntAttr)
		if intVal == 0 {
			pterm.Warning.Println("The Podman is not installed or not found in default paths")
		} else if intVal <= 2 && intVal2 < 8 {
			pterm.Warning.Printf("The Podman version %d.%d is too old, minimal version is 1.8.x\n", intVal, intVal2)
		} else {
			pterm.Success.Println("Podman: OK")
		}
	}

	intVal = reportIntAttr(allAttrs, schema.AggregionProxy, reports.NetProbeAccessibleIntAttr)
	if intVal == 0 {
		pterm.Warning.Printf("Aggregion Proxy: %s is not accessible or host is not response\n", reportDetails[schema.AggregionProxy])
	} else {
		pterm.Success.Println("Aggregion Proxy: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.DockerRegistry, reports.HTTPStatusIntAttr)
	if intVal == 0 {
		pterm.Warning.Printf("Aggregion Registry: %s is not accessible or host is not response\n", reportDetails[schema.DockerRegistry])
	} else if intVal != 200 {
		pterm.Warning.Printf("Aggregion Registry: %s response with invalid status code %d. Is registry/proxy works?\n", reportDetails[schema.DockerRegistry], intVal)
	} else {
		pterm.Success.Println("Aggregion Registry: OK")
	}
}

var commonBlockchainInspecionWasDone = false

// commonBlockchainInspection .
func commonBlockchainInspection(log *logrus.Entry, limits schema.ResourceLimitsType, allAttrs map[string]interface{}, reportDetails map[string]string) {
	if commonBlockchainInspecionWasDone {
		return
	}

	commonBlockchainInspecionWasDone = true

	var intVal int64

	pterm.DefaultSection.Println("Blockchain")

	intVal = reportIntAttr(allAttrs, schema.EOSTestNet, reports.HTTPStatusIntAttr)
	if intVal != 200 {
		pterm.Warning.Printf("TestNET: %s EOS production network is not accessible or host is not response\n", reportDetails[schema.EOSTestNet])
	} else {
		pterm.Success.Println("TestNET: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.EOSProdNet, reports.HTTPStatusIntAttr)
	if intVal != 200 {
		pterm.Warning.Printf("ProdNET: %s EOS production network is not accessible or host is not response\n", reportDetails[schema.EOSProdNet])
	} else {
		pterm.Success.Println("ProdNET: OK")
	}
}
