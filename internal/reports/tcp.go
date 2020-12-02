package reports

import "context"

const (
	// TCPAccessibleIntAttr .
	TCPAccessibleIntAttr = "accessible"
)

type (
	// TCPReport .
	TCPReport struct {
		accessible string `attr:"accessible"`
	}
)

func (dr *TCPReport) gatherLinux(ctx context.Context) error {
	dr.accessible = ""

	return nil
}

// Gather .
func (dr *TCPReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *TCPReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *TCPReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
