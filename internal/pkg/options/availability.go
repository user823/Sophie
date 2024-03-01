package options

import flag "github.com/spf13/pflag"

type AvailabilityOptions struct {
	TraceEnable   bool   `json:"trace_enable" mapstructure:"trace_enable"`
	MetricEnable  bool   `json:"metric_enable" mapstructure:"metric_enable"`
	TraceEndpoint string `json:"trace_endpoint" mapstructure:"trace_endpoint"`
	// 开启性能监控
	EnableProfiling bool `json:"profiling" mapstructure:"profiling"`
}

func NewAvailabilityOptions() *AvailabilityOptions {
	return &AvailabilityOptions{
		TraceEnable:     false,
		MetricEnable:    false,
		TraceEndpoint:   "127.0.0.1:4317",
		EnableProfiling: false,
	}
}

func (o *AvailabilityOptions) Validate() error { return nil }

func (o *AvailabilityOptions) AddFlags(fs *flag.FlagSet) {
	flag.BoolVar(&o.TraceEnable, "tracing.trace_enable", o.TraceEnable, ""+
		"Enable server to report invoke information.")
	flag.BoolVar(&o.MetricEnable, "metric.metric_enable", o.MetricEnable, ""+
		"Enable server to report source usage information.")
	flag.StringVar(&o.TraceEndpoint, "tracing.trace_endpoint", o.TraceEndpoint, ""+
		"Set server trace information report address.")
	fs.BoolVar(&o.EnableProfiling, "server.enable_profiling", o.EnableProfiling, ""+
		"Enable profiling via web interface host:port/debug/pprof/")
}
