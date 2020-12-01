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
		Total int64
		Free  int64
	}
)

func (dr *DiskReport) gatherLinux(ctx context.Context) error {
	dr.Total = 0
	dr.Free = 0

	for _, mnt := range []string{"/", "/home", "/usr", "/var", "/opt", "/aggregion"} {
		line, err := getOutputAndRegexpFind(ctx, `([\d,]+.\s+[\d,]+.\s+[\d,]+.\s+[\d]+%)\s`+mnt+"$", "df", "-m", mnt)
		if err == nil && len(line) > 0 {
			var total int64
			var use int64
			var free int64
			fmt.Sscanf(line, "%d %d %d", &total, &use, &free)
			dr.Total += total
			dr.Free += free
		}
	}

	return nil
}

// Gather .
func (dr *DiskReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *DiskReport) GetInt64(attrName string) int64 {
	switch attrName {
	case DiskTotalSpaceIntAttr:
		return dr.Total
	case DiskFreeSpaceIntAttr:
		return dr.Free
	default:
	}
	return 0
}

// GetString .
func (dr *DiskReport) GetString(attrName string) string {
	return ""
}
