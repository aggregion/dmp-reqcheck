package reports

import "context"

type (
	// IReport .
	IReport interface {
		Gather(context.Context) error
		GetInt64(attrName string) int64
		GetString(attrName string) string
	}
)
