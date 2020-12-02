package reports

import "context"

const (
	// HTTPStatusIntAttr .
	HTTPStatusIntAttr = "status"
	// HTTPHeaderStrAttrPrefix .
	HTTPHeaderStrAttrPrefix = "header_"
	// HTTPBodyStrAttr .
	HTTPBodyStrAttr = "body"
)

type (
	// HTTPReport .
	HTTPReport struct {
		status  string                 `attr:"status"`
		headers map[string]interface{} `attrMap:"header_"`
		body    string                 `attr:"body"`
	}
)

func (dr *HTTPReport) gatherLinux(ctx context.Context) error {
	dr.status = ""
	dr.headers = make(map[string]interface{})
	dr.body = ""

	return nil
}

// Gather .
func (dr *HTTPReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *HTTPReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *HTTPReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
