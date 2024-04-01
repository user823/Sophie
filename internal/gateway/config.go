package gateway

import (
	"crypto/tls"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/pkg/db/kv"
	"net"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	hserver "github.com/cloudwego/hertz/pkg/app/server"
	config "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/hertz-contrib/binding/go_playground"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/user823/Sophie/internal/gateway/router"
	"github.com/user823/Sophie/internal/pkg/middleware"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
)

// 创建、运行服务必要的配置
type Config struct {
	// 表示是否开启https
	HttpsRequired    bool
	ServerRunOptions *options.GenericRunOptions
	SecureServing    *SecureServingInfo
	InsecureServing  *InsecureServingInfo
	Discover         *DiscoverInfo
	Jwt              *router.JwtInfo
	Redis            *kv.RedisConfig
	Middlewares      []app.HandlerFunc
	Healthz          bool
	Log              *log.Options
	Aggregation      *aggregation.AnalyticsOptions
	RPCClient        *options.RPCClientOptions
	Availability     *options.AvailabilityOptions
}

func CreateConfigFromOptions(opts *Options) (*Config, error) {
	var secureServing *SecureServingInfo
	if opts.SecureServing.Required {
		secureServing = &SecureServingInfo{
			Addr:      net.JoinHostPort(opts.SecureServing.BindAddress, strconv.Itoa(opts.SecureServing.BindPort)),
			TLSConfig: opts.SecureServing.GenerateTLSConfig(),
		}
	}
	insecureServing := &InsecureServingInfo{
		Addr: net.JoinHostPort(opts.InsecureServing.BindAddress, strconv.Itoa(opts.InsecureServing.BindPort)),
	}
	discover := &DiscoverInfo{
		Addrs:    opts.ServiceDiscover.Addrs,
		Username: opts.ServiceDiscover.Username,
		Password: opts.ServiceDiscover.Password,
	}
	jwt := &router.JwtInfo{
		Realm:      opts.Jwt.Realm,
		Key:        opts.Jwt.Key,
		Timeout:    opts.Jwt.Timeout,
		MaxRefresh: opts.Jwt.MaxRefresh,
	}
	redisConfig := opts.RedisOptions.BuildRdsConfig()

	// 需要使用的中间件
	middlewares := middleware.Get("recovery", "requestid")
	middlewares = append(middlewares, middleware.Get(opts.ServerRunOptions.Middlewares...)...)

	return &Config{
		HttpsRequired:    opts.SecureServing.Required,
		SecureServing:    secureServing,
		InsecureServing:  insecureServing,
		Discover:         discover,
		Jwt:              jwt,
		Redis:            redisConfig,
		Middlewares:      middlewares,
		Healthz:          opts.ServerRunOptions.Healthz,
		Log:              opts.Log,
		Aggregation:      opts.AggregationOptions,
		RPCClient:        opts.RPCClient,
		ServerRunOptions: opts.ServerRunOptions,
		Availability:     opts.Availability,
	}, nil
}

// 创建服务通用配置
// 服务地址、链路跟踪、SSL选项需要单独开
func (cfg *Config) CreateHertzOptions() (result []config.Option) {
	// 自动裁剪路由末尾 '/'
	result = append(result, hserver.WithRedirectTrailingSlash(true))
	result = append(result, hserver.WithRemoveExtraSlash(true))

	// 优雅关停等待时间
	result = append(result, hserver.WithExitWaitTime(cfg.ServerRunOptions.ExitWaitTime))

	// 不允许使用h2c
	result = append(result, hserver.WithH2C(false))

	// 设置空闲连接时间
	result = append(result, hserver.WithIdleTimeout(cfg.ServerRunOptions.MaxIdleTimeout))

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

func (cfg *Config) CreateRemoteClientOptions() (result []client.Option) {
	// 长连接配置
	idleConfig := connpool.IdleConfig{
		MaxIdlePerAddress: cfg.RPCClient.MaxIdlePerAddress,
		MaxIdleGlobal:     cfg.RPCClient.MaxIdleGlobal,
		MaxIdleTimeout:    cfg.RPCClient.MaxIdleTimeout,
		MinIdlePerAddress: cfg.RPCClient.MinIdlePerAddress,
	}
	result = append(result, client.WithLongConnection(idleConfig))

	// 服务发现配置
	r, err := etcd.NewEtcdResolver(cfg.Discover.Addrs,
		etcd.WithAuthOpt(cfg.Discover.Username, cfg.Discover.Password))
	if err != nil {
		log.Panic(err)
	}
	result = append(result, client.WithResolver(r))

	// 超时控制
	result = append(result, client.WithConnectTimeout(cfg.RPCClient.ConnTimeout))
	result = append(result, client.WithRPCTimeout(cfg.RPCClient.RPCTimeout))

	// 服务熔断
	cbs := circuitbreak.NewCBSuite(circuitbreak.RPCInfo2Key)
	cbConfig := circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   cfg.RPCClient.Circuitbreak,
		MinSample: cfg.RPCClient.Minsample,
	}
	// 实例级熔断
	cbs.UpdateInstanceCBConfig(cbConfig)
	result = append(result, client.WithCircuitBreaker(cbs))

	// 异常重试
	fp := retry.NewFailurePolicy()
	fp.WithMaxRetryTimes(cfg.RPCClient.MaxRetryTimes)
	fp.WithMaxDurationMS(cfg.RPCClient.MaxDurationMS)
	fp.DisableChainRetryStop() // 关闭链路终止
	result = append(result, client.WithFailureRetry(fp))

	// 链路追踪
	result = append(result, client.WithSuite(kitextracing.NewClientSuite()))
	return
}

type SecureServingInfo struct {
	Addr      string
	TLSConfig *tls.Config
}

type InsecureServingInfo struct {
	Addr string
}

type DiscoverInfo struct {
	Addrs    []string
	Username string
	Password string
}
