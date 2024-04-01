package options

import (
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/errors"
)

type RMQOptions struct {
	Endpoints    []string `json:"rmq_endpoints" mapstructure:"rmq_endpoints"`
	AccessKey    string   `json:"access_key" mapstructure:"access_key"`
	AccessSecret string   `json:"access_secret" mapstructure:"access_secret"`
}

func (o *RMQOptions) Validate() error {
	if len(o.Endpoints) <= 0 {
		return errors.New("rocketmq must has at least one endpoint")
	}

	return nil
}

func (o *RMQOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringSliceVar(&o.Endpoints, "sub_rocketmq.rmq_endpoints", o.Endpoints, ""+
		"set rocketmq consumer endpoints")
	fs.StringVar(&o.AccessKey, "sub_rocketmq.access_key", o.AccessKey, ""+
		"set rocketmq consumer accessKey")
	fs.StringVar(&o.AccessSecret, "sub_rocketmq.access_secret", o.AccessSecret, ""+
		"set rocketmq consumer accessSecret")
}
