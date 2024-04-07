package aggregation

import flag "github.com/spf13/pflag"

type RMQProducerOptions struct {
	Endpoints    string `json:"endpoints" mapstructure:"endpoints"`
	AccessKey    string `json:"access_key" mapstructure:"access_key"`
	AccessSecret string `json:"access_secret" mapstructure:"access_secret"`
}

func (o *RMQProducerOptions) Validate() error { return nil }

func (o *RMQProducerOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Endpoints, "rocketmq.endpoints", o.Endpoints, "set rocketmq endpoints")
	fs.StringVar(&o.AccessKey, "rocketmq.access_key", o.AccessKey, "set rocketmq accessKey")
	fs.StringVar(&o.AccessSecret, "rocketmq.access_secret", o.AccessSecret, "set rocketmq accessSecret")
}
