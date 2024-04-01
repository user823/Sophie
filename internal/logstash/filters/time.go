package filters

import (
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

const (
	TimeFilter = "time_filter"
)

type timeFilter struct {
	cond Condition
}

// 返回一个过滤当前时间之前消息的时间过滤器
func NewTimeFilter() RecordFilter {
	created := time.Now()
	cond := func(record aggregation.LogRecord) bool {
		parsedTime := utils.Str2Time(record.Time)
		return parsedTime.After(created)
	}
	return &timeFilter{cond}
}

func (f *timeFilter) Filter(records []aggregation.LogRecord) []aggregation.LogRecord {
	return defaultFilterWithCondition(records, f.cond)
}

func (f *timeFilter) FilterWithCondition(records []aggregation.LogRecord, cond Condition) []aggregation.LogRecord {
	return defaultFilterWithCondition(records, cond)
}

func (f *timeFilter) SetCondition(cond Condition) {
	f.cond = cond
}
