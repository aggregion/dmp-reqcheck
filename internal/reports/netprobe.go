package reports

import (
	"context"
	"net"
	"time"
)

const (
	// NetProbeAccessibleIntAttr .
	NetProbeAccessibleIntAttr = "accessible"
	// NetProbeTimeIntAttr .
	NetProbeTimeIntAttr = "time"
)

type (
	// NetProbeReport .
	NetProbeReport struct {
		Type    string
		Target  string
		Timeout time.Duration

		accessible int64 `attr:"accessible"`
		probeTime  int64 `attr:"time"`
	}
)

func (dr *NetProbeReport) gatherLinux(ctx context.Context) error {
	dr.accessible = 0

	timeout := dr.Timeout
	if timeout == 0 {
		timeout = time.Second * 4
	}

	start := time.Now().UnixNano()

	conn, err := net.DialTimeout(dr.Type, dr.Target, timeout)
	if conn != nil {
		defer conn.Close()
	}
	if err == nil {
		dr.accessible = 1
	}

	dr.probeTime = (time.Now().UnixNano() - start) / 1000000

	return nil
}

// Gather .
func (dr *NetProbeReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *NetProbeReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *NetProbeReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
