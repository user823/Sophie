package options

import (
	"fmt"
	hserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/binding/go_playground"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/internal/logstash/exporters"
	"github.com/user823/Sophie/internal/logstash/filters"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/app"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/log"
	"time"
)

type Options struct {
	ServerRunOptions *options.GenericRunOptions      `json:"generic" mapstructure:"generic"`
	InsecureServing  *options.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	Log              *log.Options                    `json:"log" mapstructure:"log"`

	// 日志订阅端（目前只支持两个, 必须指定默认为空）
	SubRedisOptions *options.RedisOptions `json:"sub_redis" mapstructure:"sub_redis"`
	SubRMQOptions   *RMQOptions           `json:"sub_rocketmq" mapstructure:"sub_rocketmq"`

	// 日志输出端（默认使用stdoutExporter)
	FilterOptions *FilterOptions     `json:"filters" mapstructure:"filters"`
	Timeout       time.Duration      `json:"timeout" mapstructure:"timeout"`
	PubESOptions  *ESExporterOptions `json:"pub_elasticsearch" mapstructure:"pub_elasticsearch"`
	Exporters     []string           `json:"exporters" mapstructure:"exporters"`
}

func NewOptions() *Options {
	return &Options{
		ServerRunOptions: options.NewGenericRunOptions(),
		InsecureServing:  options.NewInsecureServingOptions(),
		Log:              log.DefaultOptions(),
		Timeout:          3 * time.Second,
		FilterOptions:    NewFilterOptions(),
		Exporters:        []string{exporters.EmptyExporter},
	}
}

func (o *Options) Flags() *ds.FlagGroup {
	fss := ds.NewFlagGroup()
	o.ServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.FilterOptions.AddFlags(fss.FlagSet("filters"))
	fss.FlagSet("common").DurationVar(&o.Timeout, "timeout", o.Timeout, "set export action timeout")
	fss.FlagSet("common").StringSliceVar(&o.Exporters, "exporters", o.Exporters, fmt.Sprintf(""+
		"set expoters, supported %v", exporters.SupportedExporters))
	return fss
}

func (o *Options) Validate() error {
	fields := []*struct {
		field interface{}
		name  string
	}{
		{o.ServerRunOptions, "generic"},
		{o.InsecureServing, "insecure"},
		{o.Log, "log"},
		{o.SubRedisOptions, "sub_redis"},
		{o.SubRMQOptions, "es_redis"},
		{o.PubESOptions, "pub_redis"},
		{o.FilterOptions, "filters"},
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

func (o *Options) CreateHertzOptions() (result []config.Option) {
	// 自动裁剪路由末尾 '/'
	result = append(result, hserver.WithRedirectTrailingSlash(true))
	result = append(result, hserver.WithRemoveExtraSlash(true))

	// 优雅关停等待时间
	result = append(result, hserver.WithExitWaitTime(o.ServerRunOptions.ExitWaitTime))

	// 不允许使用h2c
	result = append(result, hserver.WithH2C(false))

	// 设置空闲连接时间
	result = append(result, hserver.WithIdleTimeout(o.ServerRunOptions.MaxIdleTimeout))

	// 设置最大请求体(5 MB)
	result = append(result, hserver.WithMaxKeepBodySize(5*1024))

	// 自定义json 绑定器、验证器
	bindConfig := binding.NewBindConfig()
	bindConfig.UseThirdPartyJSONUnmarshaler(func(data []byte, v interface{}) error {
		return jsoniter.Unmarshal(data, v)
	})
	vd := go_playground.NewValidator()
	vd.SetValidateTag("vd")
	bindConfig.Validator = vd
	result = append(result, hserver.WithBindConfig(bindConfig))

	return
}

func (opts *Options) InitExporterConfigs() {
	fs := make([]filters.RecordFilter, len(filters.SupportedFilters))
	cnt := 0
	if opts.FilterOptions.Level >= 0 {
		fs[cnt] = filters.NewLevelFilter(log.Level(opts.FilterOptions.Level))
		cnt++
	}
	if opts.FilterOptions.TimeFilterEnable {
		fs[cnt] = filters.NewTimeFilter()
		cnt++
	}
	if len(opts.FilterOptions.LoggerWhiteList) > 0 || len(opts.FilterOptions.LoggerBlackList) > 0 {
		fs[cnt] = filters.NewLoggerNameFilter(opts.FilterOptions.LoggerBlackList, opts.FilterOptions.LoggerWhiteList)
		cnt++
	}
	fs = fs[:cnt]
	filter := fs[0]
	if cnt > 1 {
		filter = filters.NewCompositeFilter(fs...)
	}

	commonConfig := exporters.CommonConfig{}
	commonConfig.SetTimeout(opts.Timeout)
	commonConfig.SetFilter(filter)

	// 初始化StdoutExporter
	exporters.AvailableExporterConfigs.Store(exporters.StdoutExporter, &exporters.StdoutExporterConfig{
		CommonConfig: commonConfig,
	})

	// 初始化ESExporter
	if opts.PubESOptions != nil {
		esConfig := opts.PubESOptions.ESOptions.BuildESConfig()
		exporters.AvailableExporterConfigs.Store(exporters.ElasticsearchExporter, &exporters.ElasticsearchExporterConfig{
			CommonConfig: commonConfig,
			Targets:      opts.PubESOptions.TargetIndices,
			ESConfig:     esConfig,
		})
	}
}
