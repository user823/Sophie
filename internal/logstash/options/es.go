package options

import (
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

type ESExporterOptions struct {
	ESOptions     *options.ESOptions `json:"elasticsearch" mapstructure:"elasticsearch"`
	TargetIndices []string           `json:"target_indices" mapstructure:"target_indices"`
}

func (o *ESExporterOptions) Validate() error {
	if err := o.ESOptions.Validate(); err != nil {
		return err
	}

	if len(o.TargetIndices) == 0 {
		return errors.New("ESExporter at least has one target index")
	}

	return nil
}

func (o *ESExporterOptions) AddFlags(fs *flag.FlagSet) {
	o.ESOptions.AddFlags(fs)
	fs.StringSliceVar(&o.TargetIndices, "elasticsearch_exporter.target_indices", o.TargetIndices, ""+
		"set elasticsearch exporter target indices")
}

func (o *ESExporterOptions) Complete() error {
	if len(o.TargetIndices) == 0 {
		o.TargetIndices = []string{aggregation.RecordAggregation}
	}
	return nil
}
