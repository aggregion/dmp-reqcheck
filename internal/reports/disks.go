package reports

import (
	"context"
	"fmt"
	"regexp"
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

	dfOut, err := getOutput(ctx, "df", "-mT")
	if err == nil {
		rexp := regexp.MustCompile(`([/\w]+.\s+[\d,].+\s+[\d,]+.\s+[\d,]+\s+[\d]+%)\s`)

		for _, line := range rexp.FindAllString(dfOut, -1) {
			var total int64
			var use int64
			var free int64
			var fsType string

			fmt.Sscanf(line, "%s %d %d %d", &fsType, &total, &use, &free)

			switch fsType {
			case "brtfs":
			case "zfs":
			case "xfs":
			case "ext3":
			case "ext4":
				break
			default:
				continue
			}

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
