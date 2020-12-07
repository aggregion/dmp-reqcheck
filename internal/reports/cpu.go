package reports

// #cgo CFLAGS: -g -Wall
// #include "cpuid_x86.h"
import "C"
import (
	"context"
)

const (
	// CPUCoresIntAttr .
	CPUCoresIntAttr = "cores"
	// CPUFreqIntAttr .
	CPUFreqIntAttr = "freq"
	// CPUVendorStrAttr .
	CPUVendorStrAttr = "vendor"
	// CPUSmxIntAttr .
	CPUSmxIntAttr = "smx"
	// CPUSgxIntAttr .
	CPUSgxIntAttr = "sgx"
	// CPUSgx1IntAttr .
	CPUSgx1IntAttr = "sgx1"
	// CPUSgx2IntAttr .
	CPUSgx2IntAttr = "sgx2"
	// CPUSgxFlcIntAttr .
	CPUSgxFlcIntAttr = "sgx_flc"
)

type (
	// CPUReport .
	CPUReport struct {
		cores  int64  `attr:"cores"`
		freq   int64  `attr:"freq"`
		vendor string `attr:"vendor"`

		smx    int64 `attr:"smx"`
		sgx    int64 `attr:"sgx"`
		sgx1   int64 `attr:"sgx1"`
		sgx2   int64 `attr:"sgx2"`
		sgxFlc int64 `attr:"sgx_flc"`
	}

	sgxInfo struct {
		Smx       bool
		Available bool
		Version1  bool
		Version2  bool
		Flc       bool
	}
)

var _ = (IReport)((*CPUReport)(nil))

// Start .
func (dr *CPUReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *CPUReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *CPUReport) gatherLinux(ctx context.Context) []error {
	sysInfo := getSysInfo()

	dr.cores = int64(sysInfo.CPU.Cores)

	if dr.cores == 0 {
		dr.cores = int64(sysInfo.CPU.Threads)
	}

	dr.freq = int64(sysInfo.CPU.Speed)
	dr.vendor = sysInfo.CPU.Vendor

	dr.smx = 0
	dr.sgx1 = 0
	dr.sgx2 = 0
	dr.sgxFlc = 0

	sgxInfo, _ := getSgxInfo()

	if sgxInfo.Smx {
		dr.smx = 1
	}
	if sgxInfo.Available {
		dr.sgx = 1
	}
	if sgxInfo.Version1 {
		dr.sgx1 = 1
	}
	if sgxInfo.Version2 {
		dr.sgx2 = 1
	}
	if sgxInfo.Flc {
		dr.sgxFlc = 1
	}

	return nil
}

// String .
func (dr *CPUReport) String() string {
	return ""
}

// Gather .
func (dr *CPUReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *CPUReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *CPUReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
