package exporters

import (
	"github.com/user823/Sophie/pkg/errors"
	"sync"
)

var (
	AvailableExporterConfigs sync.Map
)

func NewExporter(name string) (RecordExporter, error) {
	if cfg, ok := AvailableExporterConfigs.Load(name); !ok {
		return nil, errors.Errorf("config %s is not found", name)
	} else {
		switch name {
		case StdoutExporter:
			return NewStdioExporter(cfg)
		case ElasticsearchExporter:
			return NewElasticsearchExporter(cfg)
		default:
			return nil, errors.Errorf("config %s is not found", name)
		}
	}
}
