package filters

import "github.com/user823/Sophie/pkg/log/aggregation"

type CompositeFilter struct {
	filters []RecordFilter
}

func NewCompositeFilter(filters ...RecordFilter) RecordFilter {
	return &CompositeFilter{filters}
}

// 按照子filter 加入顺序进行过滤
func (f *CompositeFilter) Filter(records []aggregation.LogRecord) []aggregation.LogRecord {
	for i := range f.filters {
		records = f.filters[i].Filter(records)
	}
	return records
}

func (f *CompositeFilter) FilterWithCondition(records []aggregation.LogRecord, cond Condition) []aggregation.LogRecord {
	for i := range f.filters {
		records = f.filters[i].FilterWithCondition(records, cond)
	}
	return records
}

// 递归设置每个子filter 的条件
func (f *CompositeFilter) SetCondition(cond Condition) {
	for i := range f.filters {
		f.filters[i].SetCondition(cond)
	}
}
