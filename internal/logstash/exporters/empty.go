package exporters

import (
	"context"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

const (
	EmptyExporter = "elasticsearch_exporter"
)

type emptyExporter struct {
	CommonConfig
}

func NewEmptyExporter() RecordExporter {
	return &emptyExporter{}
}

func (e *emptyExporter) Name() string {
	return EmptyExporter
}

func (e *emptyExporter) WriteData(context.Context, []aggregation.LogRecord) error {
	return nil
}
