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
	// CPUSgxFclIntAttr .
	CPUSgxFclIntAttr = "sgx_fcl"
)

type (
	// CPUReport .
	CPUReport struct {
		Cores  int64  `attr:"cores"`
		Freq   int64  `attr:"freq"`
		Vendor string `attr:"vendor"`

		Smx    int64 `attr:"smx"`
		Sgx    int64 `attr:"sgx"`
		Sgx1   int64 `attr:"sgx1"`
		Sgx2   int64 `attr:"sgx2"`
		SgxFcl int64 `attr:"sgx_fcl"`
	}

	sgxInfo struct {
		Smx       bool
		Available bool
		Version1  bool
		Version2  bool
		Fcl       bool
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

	dr.Cores = int64(sysInfo.CPU.Cores)

	if dr.Cores == 0 {
		dr.Cores = int64(sysInfo.CPU.Threads)
	}

	dr.Freq = int64(sysInfo.CPU.Speed)
	dr.Vendor = sysInfo.CPU.Vendor

	dr.Smx = 0
	dr.Sgx1 = 0
	dr.Sgx2 = 0
	dr.SgxFcl = 0

	sgxInfo, _ := getSgxInfo()

	if sgxInfo.Smx {
		dr.Smx = 1
	}
	if sgxInfo.Available {
		dr.Sgx = 1
	}
	if sgxInfo.Version1 {
		dr.Sgx1 = 1
	}
	if sgxInfo.Version2 {
		dr.Sgx2 = 1
	}
	if sgxInfo.Fcl {
		dr.SgxFcl = 1
	}

	return nil
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
