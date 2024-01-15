package options

import flag "github.com/spf13/pflag"

type GenericRunOptions struct {
	// 健康检查开启标志
	Healthz bool `json:"healthz" mapstructure:"healthz"`
	// 需要注册的中间件
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewGenericRunOptions() *GenericRunOptions {
	return &GenericRunOptions{
		Healthz:     false,
		Middlewares: []string{},
	}
}

func (o *GenericRunOptions) Validate() error { return nil }

func (o *GenericRunOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.BoolVar(&o.Healthz, "server.healthz", o.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&o.Middlewares, "server.middlewares", o.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.")
}
