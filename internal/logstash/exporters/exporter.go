package exporters

import (
	"context"
	"github.com/user823/Sophie/internal/logstash/filters"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"time"
)

var (
	SupportedExporters = []string{EmptyExporter, StdoutExporter, ElasticsearchExporter}
)

type RecordExporter interface {
	Name() string
	WriteData(context.Context, []aggregation.LogRecord) error
	SetFilter(filter filters.RecordFilter)
	SetTimeout(timeout time.Duration)
}
