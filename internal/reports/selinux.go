package reports

import (
	"context"
	"strings"
)

const (
	// OSSeLinuxModeStrAttr .
	OSSeLinuxModeStrAttr = "mode"
)

type (
	// OSSeLinuxReport .
	OSSeLinuxReport struct {
		mode string `attr:"mode"`

		errors []error
	}
)

var _ = (IReport)((*OSSeLinuxReport)(nil))

// Start .
func (dr *OSSeLinuxReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *OSSeLinuxReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *OSSeLinuxReport) gatherLinux(ctx context.Context) []error {
	dr.mode = ""

	output, err := getOutputAndRegexpFind(ctx, `SELinux\s+status:\s+.+`, "sestatus")
	if err == nil {
		outputParts := strings.Split(output, ":")
		if len(outputParts) == 2 {
			dr.mode = strings.TrimSpace(outputParts[1])
		}
	}

	return nil
}

// Gather .
func (dr *OSSeLinuxReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *OSSeLinuxReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *OSSeLinuxReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
