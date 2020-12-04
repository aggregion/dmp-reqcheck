package reports

import (
	"context"
)

const (
	// KernelArchStrAttr .
	KernelArchStrAttr = "arch"
	// KernelVersionStrAttr .
	KernelVersionStrAttr = "version"
	// KernelMajorVersionIntAttr .
	KernelMajorVersionIntAttr = "version_major"
	// KernelMinorVersionIntAttr .
	KernelMinorVersionIntAttr = "version_minor"
)

type (
	// KernelReport .
	KernelReport struct {
		arch         string `attr:"arch"`
		version      string `attr:"version"`
		versionMajor int64  `attr:"version_major"`
		versionMinor int64  `attr:"version_minor"`
	}
)

var _ = (IReport)((*KernelReport)(nil))

// Start .
func (dr *KernelReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *KernelReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *KernelReport) gatherLinux(ctx context.Context) []error {
	sysInfo := getSysInfo()

	dr.arch = sysInfo.Kernel.Architecture
	dr.version = sysInfo.Kernel.Release
	dr.versionMajor = 0
	dr.versionMinor = 0

	if len(dr.version) > 0 {
		dr.versionMajor, dr.versionMinor = parseVersionMinorMajor(dr.version)
	}

	return nil
}

// Gather .
func (dr *KernelReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// String .
func (dr *KernelReport) String() string {
	return ""
}

// GetInt64 .
func (dr *KernelReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *KernelReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
