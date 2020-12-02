package reports

import (
	"context"
	"fmt"
)

const (
	// RAMTotalIntAttr .
	RAMTotalIntAttr = "total"
	// RAMFreeIntAttr .
	RAMFreeIntAttr = "free"
)

type (
	// RAMReport .
	RAMReport struct {
		total int64 `attr:"total"`
		free  int64 `attr:"free"`
	}
)

func (dr *RAMReport) gatherLinux(ctx context.Context) error {
	sysInfo := getSysInfo()

	dr.total = int64(sysInfo.Memory.Size)
	dr.free = 0

	output, err := getOutputAndRegexpFind(ctx, `Mem:+\s+[\d]+\s+[\d]+\s+[\d]+`, "free", "-wm")
	if err == nil {
		var total int64
		var used int64
		fmt.Sscanf(output, "Mem: %d %d %d", &total, &used, &dr.free)
		if dr.total == 0 {
			dr.total = total
		}
	}

	return nil
}

// Gather .
func (dr *RAMReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *RAMReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *RAMReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
