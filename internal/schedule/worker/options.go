package worker

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/app"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

type Options struct {
	Log                *log.Options                    `json:"log" mapstructure:"log"`
	RedisOptions       *options.RedisOptions           `json:"redis" mapstructure:"redis"`
	AggregationOptions *aggregation.AnalyticsOptions   `json:"aggregation" mapstructure:"aggregation"`
	Availability       *options.AvailabilityOptions    `json:"availability" mapstructure:"availability"`
	RPCServer          *options.RPCServerOptions       `json:"rpc_server" mapstructure:"rpc_server"`
	MySQLOptions       *options.MySQLOptions           `json:"mysql" mapstructure:"mysql"`
	ServiceRegister    *options.ServiceDiscoverOptions `json:"service_register" mapstructure:"service_register"`
}

func NewOptions() *Options {
	return &Options{
		Log:                log.DefaultOptions(),
		RPCServer:          options.NewRPCServerOptions(),
		MySQLOptions:       options.NewMySQLOptions(),
		AggregationOptions: aggregation.NewAnalyticsOptions(),
		Availability:       options.NewAvailabilityOptions(),
		RedisOptions:       options.NewRedisOptions(),
		ServiceRegister:    options.NewServiceDiscoverOptions(),
	}
}

func (o *Options) Flags() *ds.FlagGroup {
	fss := ds.NewFlagGroup()
	o.Log.AddFlags(fss.FlagSet("log"))
	o.RPCServer.AddFlags(fss.FlagSet("remoting server"))
	o.ServiceRegister.AddFlags(fss.FlagSet("service register"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
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
		{o.MySQLOptions, "mysql"},
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
