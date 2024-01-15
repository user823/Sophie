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
		ObserverDelay:  20 * time.Second,
		RetryDelay:     10 * time.Second,
	}
}

func (o *ServiceDiscoverOptions) Validate() error { return nil }

func (o *ServiceDiscoverOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringSliceVar(&o.Addrs, "etcd.addrs", o.Addrs, "A set of etcd address(format: 127.0.0.1:2379).")
	fs.StringVar(&o.Username, "etcd.username", o.Username, "Username used to login etcd. ")
	fs.StringVar(&o.Password, "etcd.password", o.Password, "Password used to login etcd. ")
	fs.IntVar(&o.MaxAttemtTimes, "etcd.max_retry_times", o.MaxAttemtTimes, "The number of etcd connect retry times. ")
	fs.DurationVar(&o.ObserverDelay, "etcd.observer_delay", o.ObserverDelay, "ObserverDelay used to connect etcd. ")
	fs.DurationVar(&o.RetryDelay, "etcd.retry_delay", o.RetryDelay, "RetryDelay used to re-connect etcd. ")
}
