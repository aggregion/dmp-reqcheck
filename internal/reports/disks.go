package reports

import (
	"context"
	"fmt"
)

const (
	// DisksTotalSpaceIntAttr .
	DisksTotalSpaceIntAttr = "total"
	// DisksFreeSpaceIntAttr .
	DisksFreeSpaceIntAttr = "free"
)

type (
	// DisksReport .
	DisksReport struct {
		total int64 `attr:"total"`
		free  int64 `attr:"free"`
	}
)

var _ = (IReport)((*DisksReport)(nil))

// Start .
func (dr *DisksReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *DisksReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *DisksReport) gatherLinux(ctx context.Context) []error {
	dr.total = 0
	dr.free = 0

	for _, mnt := range []string{"/", "/home", "/usr", "/var", "/opt", "/aggregion", "/mount", "/mnt", "/data"} {
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

// String .
func (dr *DisksReport) String() string {
	return ""
}

// Gather .
func (dr *DisksReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *DisksReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *DisksReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
