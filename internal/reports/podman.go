package reports

import (
	"context"
	"os"
	"path"
	"time"
)

const (
	// PodmanVersionStrAttr .
	PodmanVersionStrAttr = "version"
	// PodmanMajorVersionIntAttr .
	PodmanMajorVersionIntAttr = "version_major"
	// PodmanMinorVersionIntAttr .
	PodmanMinorVersionIntAttr = "version_minor"
)

type (
	// PodmanReport .
	PodmanReport struct {
		version      string `attr:"version"`
		versionMajor int64  `attr:"version_major"`
		versionMinor int64  `attr:"version_minor"`
	}
)

var _ = (IReport)((*PodmanReport)(nil))

// Start .
func (dr *PodmanReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *PodmanReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *PodmanReport) gatherLinux(ctx context.Context) []error {
	dr.version = ""
	dr.versionMajor, dr.versionMinor = 0, 0

	homeDir, err := os.UserHomeDir()
	paths := findLinuxApps(
		context.Background(),
		[]string{"/usr", "/opt", path.Join(homeDir, ".local"), path.Join(homeDir, "bin")},
		[]string{
			"podman",
		},
	)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	version, err := getOutputAndRegexpFind(ctx, versionRegExp, paths["podman"], "--version")
	if err == nil {
		dr.version = version
		dr.versionMajor, dr.versionMinor = parseVersionMinorMajor(dr.version)
	}

	return nil
}

// String .
func (dr *PodmanReport) String() string {
	return ""
}

// Gather .
func (dr *PodmanReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *PodmanReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *PodmanReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
