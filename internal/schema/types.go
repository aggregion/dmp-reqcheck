package schema

import (
	"github.com/aggregion/dmp-reqcheck/internal/config"
	"github.com/aggregion/dmp-reqcheck/internal/reports"
)

type (
	// ResourceLimitType .
	ResourceLimitType = struct {
		Minimal   int64
		Recommend int64
	}

	// ResourceLimitsType .
	ResourceLimitsType = map[string]ResourceLimitType

	// ReportsGroup .
	ReportsGroup = map[string]reports.IReport

	// CheckContext .
	CheckContext struct {
		Path  string
		Roles []string

		ResourceLimits ResourceLimitsType
		StrAttrs       map[string]string
		IntAttrs       map[string]int64
	}

	// CheckSchema .
	CheckSchema struct {
		Name           string
		Role           config.Role
		ResourceLimits ResourceLimitsType
		Reports        ReportsGroup
	}
)
