package reports

import "context"

const (
	// HypervisorNameStrAttr .
	HypervisorNameStrAttr = "name"
)

type (
	// HypervisorReport .
	HypervisorReport struct {
		name string `attr:"name"`
	}
)

func (dr *HypervisorReport) gatherLinux(ctx context.Context) error {
	sysInfo := getSysInfo()

	dr.name = sysInfo.Node.Hypervisor

	return nil
}

// Gather .
func (dr *HypervisorReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *HypervisorReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *HypervisorReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
