package reports

import (
	"context"
)

const (
	// OSArchStrAttr .
	OSArchStrAttr = "arch"
	// OSTypeStrAttr .
	OSTypeStrAttr = "type"
	// OSVendorStrAttr .
	OSVendorStrAttr = "vendor"
	// OSVersionStrAttr .
	OSVersionStrAttr = "version"
	// OSMajorVersionIntAttr .
	OSMajorVersionIntAttr = "version_major"
	// OSMinorVersionIntAttr .
	OSMinorVersionIntAttr = "version_minor"
)

type (
	// OSReport .
	OSReport struct {
		arch   string `attr:"arch"`
		osType string `attr:"type"`
		vendor string `attr:"vendor"`

		version      string `attr:"version"`
		versionMajor int64  `attr:"version_major"`
		versionMinor int64  `attr:"version_minor"`
	}
)

func (dr *OSReport) gatherLinux(ctx context.Context) error {
	sysInfo := getSysInfo()

	dr.osType = "linux"
	dr.arch = sysInfo.OS.Architecture
	dr.vendor = sysInfo.OS.Vendor
	dr.version = sysInfo.OS.Version
	dr.versionMajor, dr.versionMinor = parseVersionMinorMajor(dr.version)
	if dr.versionMajor == 0 {
		dr.versionMajor, dr.versionMinor = parseVersionMinorMajor(sysInfo.OS.Release)
	}
	if dr.versionMajor == 0 {
		dr.versionMajor, dr.versionMinor = parseVersionMinorMajor(sysInfo.OS.Name)
	}

	return nil
}

// Gather .
func (dr *OSReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *OSReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *OSReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
