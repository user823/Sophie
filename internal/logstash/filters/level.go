package filters

import (
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/utils"
	"go.uber.org/zap/zapcore"
)

const (
	LevelFilter = "level_filter"
)

type levelFilter struct {
	cond Condition
}

func NewLevelFilter(level log.Level) RecordFilter {
	cond := func(record aggregation.LogRecord) bool {
		var zapLevel zapcore.Level
		if err := zapLevel.UnmarshalText(utils.S2b(record.Level)); err != nil {
			log.Errorf("resolve record level error: %s", err.Error())
			return false
		}
		return zapLevel >= level
	}
	return &levelFilter{cond}
}

func (f *levelFilter) Filter(records []aggregation.LogRecord) []aggregation.LogRecord {
	return defaultFilterWithCondition(records, f.cond)
}

func (f *levelFilter) FilterWithCondition(records []aggregation.LogRecord, cond Condition) []aggregation.LogRecord {
	return defaultFilterWithCondition(records, cond)
}

func (f *levelFilter) SetCondition(cond Condition) {
	f.cond = cond
}
