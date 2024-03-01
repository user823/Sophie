package options

import (
	flag "github.com/spf13/pflag"
	"time"
)

// etcd 服务注册与发现配置
// 默认不开启认证
type ServiceDiscoverOptions struct {
	Addrs          []string      `json:"addrs" mapstructure:"addrs"`
	Username       string        `json:"username" mapstructure:"username"`
	Password       string        `json:"password" mapstructure:"password"`
	MaxAttemtTimes int           `json:"max_attemt_times" mapstructure:"max_attemt_times"`
	ObserverDelay  time.Duration `json:"observer_delay" mapstructure:"observer_delay"`
	RetryDelay     time.Duration `json:"retry_delay" mapstructure:"retry_delay"`
}

func NewServiceDiscoverOptions() *ServiceDiscoverOptions {
	return &ServiceDiscoverOptions{
		Addrs:          []string{"127.0.0.1:2379"},
		MaxAttemtTimes: 5,
		ObserverDelay:  10 * time.Second,
		RetryDelay:     5 * time.Second,
	}
}

func (o *ServiceDiscoverOptions) Validate() error { return nil }

func (o *ServiceDiscoverOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringSliceVar(&o.Addrs, "server_discover.addrs", o.Addrs, "A set of etcd address(format: 127.0.0.1:2379).")
	fs.StringVar(&o.Username, "server_discover.username", o.Username, "Username used to login etcd. ")
	fs.StringVar(&o.Password, "server_discover.password", o.Password, "Password used to login etcd. ")
	fs.IntVar(&o.MaxAttemtTimes, "server_discover.max_retry_times", o.MaxAttemtTimes, "The number of etcd connect retry times. ")
	fs.DurationVar(&o.ObserverDelay, "server_discover.observer_delay", o.ObserverDelay, "ObserverDelay used to connect etcd. ")
	fs.DurationVar(&o.RetryDelay, "server_discover.retry_delay", o.RetryDelay, "RetryDelay used to re-connect etcd. ")
}
