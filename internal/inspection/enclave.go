package inspection

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/logger"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
	"github.com/aggregion/dmp-reqcheck/internal/schema"
	"github.com/pterm/pterm"
)

// EnclaveInspection .
func EnclaveInspection(cfg *config.Settings, sc *schema.CheckSchema, reportDetails map[string]string) {
	limits := sc.ResourceLimits
	allAttrs := schema.MergeReportsAttrs(sc.Reports)

	log := logger.Get("check", "Enclave Inspection")

	// pterm.DefaultSection.Println("Enclave Inspection Report")

	commonInspection(log, limits, allAttrs, reportDetails)
	commonBlockchainInspection(log, limits, allAttrs, reportDetails)

	pterm.DefaultSection.Println("Enclave Specific Report")

	// var intVal, intVal2, minVal int64
	var intVal int64
	var flcVal int64
	var strVal string
	var strVal2 string

	//
	// Hardware
	//
	intVal = reportIntAttr(allAttrs, schema.CPU, reports.CPUSgx1IntAttr)
	if intVal != 1 {
		pterm.Error.Printf("CPU: SGX1 is not supported or it not enabled\n")
	} else {
		pterm.Success.Println("CPU SGX: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.CPU, reports.CPUSgx2IntAttr)
	if intVal != 0 {
		pterm.Success.Printf("CPU SGX2: OK\n")
	}

	flcVal = reportIntAttr(allAttrs, schema.CPU, reports.CPUSgxFlcIntAttr)
	if flcVal != 1 {
		pterm.Warning.Printf("CPU: FLC feature is not supported or it not enabled, not possible to use Intel DCAP driver\n")
	} else {
		pterm.Success.Println("CPU FLC: OK")
	}

	//
	// OS
	//
	strVal = reportStrAttr(allAttrs, schema.OS, reports.OSVendorStrAttr)
	if strVal != "ubuntu" && strVal != "centos" {
		pterm.Warning.Printf("OS Vendor: current OS is %s, the driver might be not works properly, should be use Ubuntu 18.x or CentOS 8.x\n", strVal)
	} else {
		pterm.Success.Println("OS Vendor: OK")
	}

	osVersionOk := true
	strVal = reportStrAttr(allAttrs, schema.OS, reports.OSVersionStrAttr)
	intVal = reportIntAttr(allAttrs, schema.OS, reports.OSMajorVersionIntAttr)
	switch strVal {
	case "ubuntu":
		if intVal < 18 {
			pterm.Warning.Printf("OS Version: Ubuntu major version %s is too old, minimal version is 18.x\n", strVal)
			osVersionOk = false
		}
	case "centos":
		if intVal < 8 {
			pterm.Warning.Printf("OS Version: CentOS major version %s is too old, minimal version is 8.x\n", strVal)
			osVersionOk = false
		}
	}

	if osVersionOk {
		pterm.Success.Println("OS Version: OK")
	}

	//
	// Drivers
	//
	strVal = reportStrAttr(allAttrs, schema.DriverDCAP, reports.DriverVersionStrAttr)
	strVal2 = reportStrAttr(allAttrs, schema.DriverISGX, reports.DriverVersionStrAttr)
	if strVal == "" && strVal2 == "" {
		pterm.Warning.Printf("The SGX Driver: no sgx drivers was found on the host\n")
	} else if strVal2 != "" && flcVal != 0 {
		pterm.Warning.Printf("The SGX Driver: ISGX driver was found on the host, your CPU is supported FLC, please install DCAP SGX Driver instead\n")
	} else if strVal != "" && flcVal == 0 {
		pterm.Warning.Printf("The SGX Driver: DCAP sgx driver was found on the host, your CPU is not supported FLC, please install Intel SGX Driver (not DCAP) instead\n")
	} else if strVal2 != "" {
		pterm.Success.Println("The SGX Driver: OK (ISGX)")
	} else if strVal != "" {
		pterm.Success.Println("The SGX Driver: OK (DCAP)")
	} else {
		pterm.Warning.Printf("The SGX Driver: no any SGX drivers was found on the host\n")
	}

	//
	// Hosts
	//
	intVal = reportIntAttr(allAttrs, schema.ClickhouseHostSvc, reports.HostAccessibleIntAttr)
	if intVal != 1 {
		pterm.Error.Printf("Clickhouse Host: the host %s is not accessible or not response\n", reportDetails[schema.ClickhouseHostSvc])
	} else {
		pterm.Success.Println("Clickhouse Host: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.DmpHostSvc, reports.HostAccessibleIntAttr)
	if intVal != 1 {
		pterm.Error.Printf("DMP Host: the host %s is not accessible or not response\n", reportDetails[schema.DmpHostSvc])
	} else {
		pterm.Success.Println("DMP Host: OK")
	}

	pterm.DefaultSection.Println("External for SGX")

	//
	// CAS
	//
	intVal = reportIntAttr(allAttrs, schema.CasTest, reports.NetProbeAccessibleIntAttr)
	if intVal != 1 {
		pterm.Warning.Printf("CAS Test Host: the host %s is not accessible or not response\n", reportDetails[schema.CasTest])
	} else {
		pterm.Success.Println("CAS Test Host: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.CasProd, reports.NetProbeAccessibleIntAttr)
	if intVal != 1 {
		pterm.Error.Printf("CAS Prod Host: the host %s is not accessible or not response\n", reportDetails[schema.CasProd])
	} else {
		pterm.Success.Println("CAS Prod Host: OK")
	}

	//
	// Intel
	//
	intVal = reportIntAttr(allAttrs, schema.IntelPbs, reports.HTTPStatusIntAttr)
	if intVal == 0 {
		pterm.Error.Printf("Intel SGX PBS: resource %s is not accessible\n", reportDetails[schema.IntelPbs])
	} else {
		pterm.Success.Println("Intel SGX PBS: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.IntelOcsp, reports.HTTPStatusIntAttr)
	if intVal == 0 {
		pterm.Error.Printf("Intel SGX OSCP: resource %s is not accessible\n", reportDetails[schema.IntelOcsp])
	} else {
		pterm.Success.Println("Intel SGX OSCP: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.IntelCrl, reports.HTTPStatusIntAttr)
	if intVal == 0 {
		pterm.Warning.Printf("Intel SGX CRL: resource %s is not accessible\n", reportDetails[schema.IntelCrl])
	} else {
		pterm.Success.Println("Intel SGX CRL: OK")
	}

	intVal = reportIntAttr(allAttrs, schema.IntelWl, reports.HTTPStatusIntAttr)
	if intVal == 0 {
		pterm.Error.Printf("Intel SGX WL: resource %s is not accessible\n", reportDetails[schema.IntelWl])
	} else {
		pterm.Success.Println("Intel SGX WL: OK")
	}
}
