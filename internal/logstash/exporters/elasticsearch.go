package exporters

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/user823/Sophie/pkg/db/doc"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

const (
	ElasticsearchExporter = "elasticsearch_exporter"
)

type elasticsearchExporter struct {
	es *elasticsearch.TypedClient
	CommonConfig
	targets []string
}

type ElasticsearchExporterConfig struct {
	Targets []string
	CommonConfig
	ESConfig *doc.ESConfig
}

func NewElasticsearchExporter(config any) (RecordExporter, error) {
	if cfg, ok := config.(*ElasticsearchExporterConfig); ok {
		escli, err := doc.NewES(cfg.ESConfig)
		if err != nil {
			return nil, errors.New("elasticsearchExporter config is invalid")
		}
		return &elasticsearchExporter{
			es:           escli,
			CommonConfig: cfg.CommonConfig,
			targets:      cfg.Targets,
		}, nil
	}
	return nil, errors.New("elasticsearchExporter config is invalid")
}

func (e *elasticsearchExporter) Name() string {
	return ElasticsearchExporter
}

func (e *elasticsearchExporter) WriteData(ctx context.Context, data []aggregation.LogRecord) error {
	if e.filter != nil {
		data = e.filter.Filter(data)
	}

	bk := e.es.Bulk()
	for i := range data {
		for j := range e.targets {
			if err := bk.CreateOp(types.CreateOperation{Index_: &e.targets[j]}, data[i]); err != nil {
				log.Warn("es exporter create record error: %s", err.Error())
				continue
			}
		}
	}
	resp, err := bk.Timeout(e.timeout.String()).Do(ctx)
	if err != nil || resp.Errors {
		return errors.New("es exporter write data failed")
	}
	return nil
}
