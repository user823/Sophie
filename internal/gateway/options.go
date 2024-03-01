package gateway

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/app"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

// 运行gateway options
// 实现App.option 若干接口
type Options struct {
	ServiceDiscover    *options.ServiceDiscoverOptions `json:"server_discover" mapstructure:"server_discover"`
	ServerRunOptions   *options.GenericRunOptions      `json:"generic" mapstructure:"generic"`
	InsecureServing    *options.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing      *options.SecureServingOptions   `json:"secure" mapstructure:"secure"`
	RedisOptions       *options.RedisOptions           `json:"redis" mapstructure:"redis"`
	RPCClient          *options.RPCClientOptions       `json:"rpc_client" mapstructure:"rpc_client"`
	Log                *log.Options                    `json:"log" mapstructure:"log"`
	Jwt                *options.JwtOptions             `json:"jwt" mapstructure:"jwt"`
	AggregationOptions *aggregation.AnalyticsOptions   `json:"aggregation" mapstructure:"aggregation"`
	Availability       *options.AvailabilityOptions    `json:"availability" mapstructure:"availability"`
}

func NewOptions() *Options {
	return &Options{
		ServiceDiscover:    options.NewServiceDiscoverOptions(),
		ServerRunOptions:   options.NewGenericRunOptions(),
		InsecureServing:    options.NewInsecureServingOptions(),
		SecureServing:      options.NewSecureServingOptions(),
		RedisOptions:       options.NewRedisOptions(),
		Log:                log.DefaultOptions(),
		RPCClient:          options.NewRPCClientOptions(),
		Jwt:                options.NewJwtOptions(),
		AggregationOptions: aggregation.NewAnalyticsOptions(),
		Availability:       options.NewAvailabilityOptions(),
	}
}

func (o *Options) Flags() *ds.FlagGroup {
	fss := ds.NewFlagGroup()
	o.ServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure"))
	o.SecureServing.AddFlags(fss.FlagSet("server"))
	o.ServiceDiscover.AddFlags(fss.FlagSet("server_discover"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.Jwt.AddFlags(fss.FlagSet("jwt"))
	o.RPCClient.AddFlags(fss.FlagSet("rpc_client"))
	o.AggregationOptions.AddFlags(fss.FlagSet("aggregation"))
	o.Availability.AddFlags(fss.FlagSet("availability"))
	return fss
}

func (o *Options) Validate() error {
	fields := []*struct {
		field interface{}
		name  string
	}{
		{o.ServiceDiscover, "server_discover"},
		{o.ServerRunOptions, "generic"},
		{o.InsecureServing, "insecure"},
		{o.SecureServing, "secure"},
		{o.RedisOptions, "redis"},
		{o.Log, "log"},
		{o.Jwt, "jwt"},
		{o.RPCClient, "rpc_client"},
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
