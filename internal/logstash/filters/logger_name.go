package filters

import (
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

const (
	LoggerNameFilter = "logger_name_filter"
)

type loggerNameFilter struct {
	cond Condition
}

// 白名单为空时默认放行所有logger
func NewLoggerNameFilter(blackNames []string, whiteNames []string) RecordFilter {
	cond := func(record aggregation.LogRecord) bool {
		if strutil.ContainsAny(record.LoggerName, blackNames...) {
			return false
		}

		if len(whiteNames) == 0 {
			return true
		}

		return strutil.ContainsAny(record.LoggerName, whiteNames...)
	}

	return &loggerNameFilter{cond}
}

func (f *loggerNameFilter) Filter(records []aggregation.LogRecord) []aggregation.LogRecord {
	return defaultFilterWithCondition(records, f.cond)
}

func (f *loggerNameFilter) FilterWithCondition(records []aggregation.LogRecord, cond Condition) []aggregation.LogRecord {
	return defaultFilterWithCondition(records, cond)
}

func (f *loggerNameFilter) SetCondition(cond Condition) {
	f.cond = cond
}
