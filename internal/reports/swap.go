package reports

import (
	"context"
	"fmt"
)

const (
	// SwapFileTotalIntAttr .
	SwapFileTotalIntAttr = "total"
	// SwapFileFreeIntAttr .
	SwapFileFreeIntAttr = "free"
)

type (
	// SwapFileReport .
	SwapFileReport struct {
		total int64 `attr:"total"`
		free  int64 `attr:"free"`
	}
)

func (dr *SwapFileReport) gatherLinux(ctx context.Context) error {
	dr.total = 0
	dr.free = 0

	output, err := getOutputAndRegexpFind(ctx, `Swap:+\s+[\d]+\s+[\d]+\s+[\d]+`, "free", "-wm")
	if err == nil {
		var used int64
		fmt.Sscanf(output, "Swap: %d %d %d", &dr.total, &used, &dr.free)
	}

	return nil
}

// Gather .
func (dr *SwapFileReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *SwapFileReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *SwapFileReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
