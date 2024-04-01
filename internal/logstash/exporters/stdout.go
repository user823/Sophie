package exporters

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

const (
	StdoutExporter = "stdio_exporter"
)

// 输出到标准控制台的exporter
type stdoutExporter struct {
	CommonConfig
}

type StdoutExporterConfig struct {
	CommonConfig
}

func NewStdioExporter(config any) (RecordExporter, error) {
	if cfg, ok := config.(*StdoutExporterConfig); ok {
		return &stdoutExporter{cfg.CommonConfig}, nil
	}
	return nil, errors.New("stdoutExporter config is invalid")
}

func (e *stdoutExporter) Name() string {
	return StdoutExporter
}

func (e *stdoutExporter) WriteData(ctx context.Context, data []aggregation.LogRecord) error {
	if e.timeout > 0 {
		var tctx context.Context
		tctx, cancel := context.WithTimeout(ctx, e.timeout)
		defer cancel()
		ctx = tctx
	}

	if e.filter != nil {
		data = e.filter.Filter(data)
	}
	for _, record := range data {
		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
		}
		fmt.Printf("%+v", record)
	}
	return nil
}
