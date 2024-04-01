package file

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/app"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

// 允许file options
// 实现App.option 若干接口
type Options struct {
	Log                *log.Options                    `json:"log" mapstructure:"log"`
	RPCServer          *options.RPCServerOptions       `json:"rpc_server" mapstructure:"rpc_server"`
	ServiceRegister    *options.ServiceDiscoverOptions `json:"service_register" mapstructure:"service_register"`
	RedisOptions       *options.RedisOptions           `json:"redis" mapstructure:"redis"`
	AggregationOptions *aggregation.AnalyticsOptions   `json:"aggregation" mapstructure:"aggregation"`
	Availability       *options.AvailabilityOptions    `json:"availability" mapstructure:"availability"`
}

func NewOptions() *Options {
	return &Options{
		Log:                log.DefaultOptions(),
		RPCServer:          options.NewRPCServerOptions(),
		ServiceRegister:    options.NewServiceDiscoverOptions(),
		RedisOptions:       options.NewRedisOptions(),
		AggregationOptions: aggregation.NewAnalyticsOptions(),
		Availability:       options.NewAvailabilityOptions(),
	}
}

func (o *Options) Flags() *ds.FlagGroup {
	fss := ds.NewFlagGroup()
	o.Log.AddFlags(fss.FlagSet("log"))
	o.RPCServer.AddFlags(fss.FlagSet("remoting server"))
	o.ServiceRegister.AddFlags(fss.FlagSet("service register"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.AggregationOptions.AddFlags(fss.FlagSet("log aggregation"))
	o.Availability.AddFlags(fss.FlagSet("availability"))
	return fss
}

func (o *Options) Validate() error {
	fields := []*struct {
		field interface{}
		name  string
	}{
		{o.Log, "log"},
		{o.RPCServer, "rpc_server"},
		{o.ServiceRegister, "service register"},
		{o.RedisOptions, "redis"},
		{o.AggregationOptions, "aggregation"},
		{o.Availability, "availability"},
	}
	for _, item := range fields {
		field, name := item.field, item.name
		if v, ok := field.(app.ValidatableOptions); ok {
			if err := v.Validate(); err != nil {
				log.Errorf("Error %s option: %s", name, err.Error())
				return err
			}
		}
	}
	return nil
}

func (o *Options) String() string {
	data, _ := jsoniter.Marshal(o)
	return string(data)
}
