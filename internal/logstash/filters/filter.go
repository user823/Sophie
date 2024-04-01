package filters

import "github.com/user823/Sophie/pkg/log/aggregation"

var (
	SupportedFilters = []string{LevelFilter, LoggerNameFilter, TimeFilter}
)

type Condition func(record aggregation.LogRecord) bool

type RecordFilter interface {
	// 使用默认过滤条件进行过滤
	Filter([]aggregation.LogRecord) []aggregation.LogRecord
	// 指定本次过滤使用的过滤条件
	FilterWithCondition([]aggregation.LogRecord, Condition) []aggregation.LogRecord
	// 设置过滤条件
	SetCondition(Condition)
}

func defaultFilterWithCondition(records []aggregation.LogRecord, cond Condition) []aggregation.LogRecord {
	cnt := 0
	for i := range records {
		if cond(records[i]) {
			records[cnt] = records[i]
			cnt++
		}
	}
	records = records[:cnt]
	return records
}
