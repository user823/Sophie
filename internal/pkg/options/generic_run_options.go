package options

import (
	"time"

	flag "github.com/spf13/pflag"
)

type GenericRunOptions struct {
	// 健康检查开启标志
	Healthz bool `json:"healthz" mapstructure:"healthz"`
	// 需要注册的中间件
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
	// 优雅关停的等待时间
	ExitWaitTime time.Duration `json:"exit_wait_time" mapstructure:"exit_wait_time"`
	// 客户端最大空闲连接时间
	MaxIdleTimeout time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`
	// api版本
	BaseAPI string `json:"base_api" mapstructure:"base_api"`
}

func NewGenericRunOptions() *GenericRunOptions {
	return &GenericRunOptions{
		Healthz:        false,
		Middlewares:    []string{},
		ExitWaitTime:   8 * time.Second,
		MaxIdleTimeout: 5 * time.Second,
	}
}

func (o *GenericRunOptions) Validate() error { return nil }

func (o *GenericRunOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.BoolVar(&o.Healthz, "generic.healthz", o.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&o.Middlewares, "generic.middlewares", o.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. Supported middlewares are cache, recovery, cors, requestid, accesslog.")
	fs.DurationVar(&o.ExitWaitTime, "generic.exit_wait_time", o.ExitWaitTime, ""+
		"Exit wait time for server shutdown")
	fs.DurationVar(&o.MaxIdleTimeout, "generic.idle_timeout", o.MaxIdleTimeout, ""+
		"Server max idle conn timeout")
	fs.StringVar(&o.BaseAPI, "generic.base_api", o.BaseAPI, ""+
		"Server api version (eg. dev-api, proc-api, v1, v2...)")
}
