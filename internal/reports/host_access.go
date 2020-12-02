package reports

import "context"

const (
	// HTTPAccessibleIntAttr .
	HTTPAccessibleIntAttr = "accessible"
)

type (
	// HostAccesseReport .
	HostAccesseReport struct {
		Accessible string `attr:"accessible"`
	}
)

func (dr *HostAccesseReport) gatherLinux(ctx context.Context) error {
	dr.Accessible = ""

	return nil
}

// Gather .
func (dr *HostAccesseReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *HostAccesseReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *HostAccesseReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
