package reports

import "context"

type (
	// IReport .
	IReport interface {
		Start(context.Context) error
		Stop(context.Context) error
		String() string
		Gather(context.Context) []error
		GetInt64(attrName string) int64
		GetString(attrName string) string
	}
)
