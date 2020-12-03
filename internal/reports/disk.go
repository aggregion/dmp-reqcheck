package reports

import (
	"context"
	"fmt"
)

const (
	// DiskTotalSpaceIntAttr .
	DiskTotalSpaceIntAttr = "total"
	// DiskFreeSpaceIntAttr .
	DiskFreeSpaceIntAttr = "free"
)

type (
	// DiskReport .
	DiskReport struct {
		total int64 `attr:"total"`
		free  int64 `attr:"free"`
	}
)

var _ = (IReport)((*DiskReport)(nil))

// Start .
func (dr *DiskReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *DiskReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *DiskReport) gatherLinux(ctx context.Context) []error {
	dr.total = 0
	dr.free = 0

	for _, mnt := range []string{"/", "/home", "/usr", "/var", "/opt", "/aggregion", "/mount", "/mnt"} {
		line, err := getOutputAndRegexpFind(ctx, `([\d,]+.\s+[\d,]+.\s+[\d,]+.\s+[\d]+%)\s`+mnt+"$", "df", "-m", mnt)
		if err == nil && len(line) > 0 {
			var total int64
			var use int64
			var free int64
			fmt.Sscanf(line, "%d %d %d", &total, &use, &free)
			dr.total += total
			dr.free += free
		}
	}

	return nil
}

// Gather .
func (dr *DiskReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *DiskReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *DiskReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
