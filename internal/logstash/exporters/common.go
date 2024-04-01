package exporters

import (
	"github.com/user823/Sophie/internal/logstash/filters"
	"github.com/user823/Sophie/pkg/errors"
	"time"
)

var (
	ErrTimeout = errors.New("time out")
)

type CommonConfig struct {
	filter  filters.RecordFilter
	timeout time.Duration
}

func (p *CommonConfig) SetFilter(filter filters.RecordFilter) {
	p.filter = filter
}

func (p *CommonConfig) SetTimeout(t time.Duration) {
	p.timeout = t
}
