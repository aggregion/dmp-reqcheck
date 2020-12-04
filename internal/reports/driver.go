package reports

import (
	"context"
	"regexp"
)

const (
	// DriverVersionStrAttr .
	DriverVersionStrAttr = "version"
	// DriverMajorVersionIntAttr .
	DriverMajorVersionIntAttr = "version_major"
	// DriverMinorVersionIntAttr .
	DriverMinorVersionIntAttr = "version_minor"
)

type (
	// DriverReport .
	DriverReport struct {
		DriverName string

		version      string `attr:"version"`
		versionMajor int64  `attr:"version_major"`
		versionMinor int64  `attr:"version_minor"`
	}
)

var _ = (IReport)((*DriverReport)(nil))

// Start .
func (dr *DriverReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *DriverReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *DriverReport) gatherLinux(ctx context.Context) []error {
	dr.version = ""
	dr.versionMajor = 0
	dr.versionMinor = 0

	version, err := getOutputAndRegexpFind(ctx, `version.+`, "modinfo", dr.DriverName)
	if err == nil {
		dr.version = regexp.MustCompile(versionRegExp).FindString(version)
		dr.versionMajor, dr.versionMinor = parseVersionMinorMajor(dr.version)
	}

	return nil
}

// String .
func (dr *DriverReport) String() string {
	return dr.DriverName
}

// Gather .
func (dr *DriverReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *DriverReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *DriverReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
