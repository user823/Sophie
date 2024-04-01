package gen

import (
	jsoniter "github.com/json-iterator/go"
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/internal/gen/engine"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

// 允许system options
// 实现App.option 若干接口
type Options struct {
	Log                *log.Options                    `json:"log" mapstructure:"log"`
	RPCServer          *options.RPCServerOptions       `json:"rpc_server" mapstructure:"rpc_server"`
	MySQLOptions       *options.MySQLOptions           `json:"mysql" mapstructure:"mysql"`
	ServiceRegister    *options.ServiceDiscoverOptions `json:"service_register" mapstructure:"service_register"`
	RedisOptions       *options.RedisOptions           `json:"redis" mapstructure:"redis"`
	AggregationOptions *aggregation.AnalyticsOptions   `json:"aggregation" mapstructure:"aggregation"`
	Availability       *options.AvailabilityOptions    `json:"availability" mapstructure:"availability"`
	GenOptions         *GenOptions                     `json:"gen" mapstructure:"gen"`
}

type GenOptions struct {
	// 引擎名称
	EngineName string `json:"engine_name" mapstructure:"engine_name"`
	// 模板搜索路径
	SearchPath []string                 `json:"template_search_path" mapstructure:"template_search_path"`
	GenParams  *engine.GenHelperOptions `json:"params" mapstructure:"params"`
}

func (o *GenOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.EngineName, "gen.engine_name", o.EngineName, ""+
		"Set gen code engine name")
	fs.StringSliceVar(&o.SearchPath, "gen.template_search_path", o.SearchPath, ""+
		"Set gen code engine search templates path (dir or file path, if dir included children path also)")
	// 递归调用
	o.GenParams.AddFlags(fs)
}

func NewOptions() *Options {
	return &Options{
		Log:                log.DefaultOptions(),
		RPCServer:          options.NewRPCServerOptions(),
		MySQLOptions:       options.NewMySQLOptions(),
		ServiceRegister:    options.NewServiceDiscoverOptions(),
		RedisOptions:       options.NewRedisOptions(),
		AggregationOptions: aggregation.NewAnalyticsOptions(),
		Availability:       options.NewAvailabilityOptions(),
		GenOptions: &GenOptions{
			EngineName: "gen code engine",
			SearchPath: []string{"."},
			GenParams:  &engine.GlobalOptions,
		},
	}
}

func (o *Options) Flags() *ds.FlagGroup {
	fss := ds.NewFlagGroup()
	o.GenOptions.AddFlags(fss.FlagSet("gen"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.RPCServer.AddFlags(fss.FlagSet("remoting server"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.ServiceRegister.AddFlags(fss.FlagSet("service register"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.AggregationOptions.AddFlags(fss.FlagSet("log aggregation"))
	o.Availability.AddFlags(fss.FlagSet("availability"))
	return fss
}

func (o *Options) String() string {
	data, _ := jsoniter.Marshal(o)
	return string(data)
}
