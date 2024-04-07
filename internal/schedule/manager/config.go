package manager

import (
	"net"
	"strconv"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	kretry "github.com/kitex-contrib/registry-etcd/retry"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/internal/schedule/manager/loadbalance"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/db/sql"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/log/aggregation/producer"
)

type Config struct {
	Log             *log.Options
	ServerRunConfig *options.RPCServerOptions
	Register        *options.ServiceDiscoverOptions
	Aggregation     *aggregation.AnalyticsOptions
	Redis           *kv.RedisConfig
	MySQL           *sql.MysqlConfig
	Availability    *options.AvailabilityOptions
	LoadBalance     *loadbalance.LoadBalanceOptions
	RPCClient       *options.RPCClientOptions
}

func CreateConfigFromOptions(opts *Options) (*Config, error) {
	// redis
	redisConfig := &kv.RedisConfig{
		Addrs:                 opts.RedisOptions.Addrs,
		MasterName:            opts.RedisOptions.MasterName,
		Username:              opts.RedisOptions.Username,
		Password:              opts.RedisOptions.Password,
		Database:              opts.RedisOptions.Database,
		MaxIdle:               opts.RedisOptions.MaxIdle,
		MaxActive:             opts.RedisOptions.MaxActive,
		Timeout:               opts.RedisOptions.Timeout,
		EnableCluster:         opts.RedisOptions.EnableCluster,
		UseSSL:                opts.RedisOptions.UseSSL,
		SSLInsecureSkipVerify: opts.RedisOptions.SSLInsecureSkipVerify,
	}

	// mysql config
	mysqlConfig := &sql.MysqlConfig{
		Host:                  opts.MySQLOptions.Host,
		Username:              opts.MySQLOptions.Username,
		Password:              opts.MySQLOptions.Password,
		Database:              opts.MySQLOptions.Database,
		MaxIdleConnections:    opts.MySQLOptions.MaxIdleConnections,
		MaxOpenConnections:    opts.MySQLOptions.MaxOpenConnections,
		MaxConnectionLifeTime: opts.MySQLOptions.MaxConnectionLifeTime,
		LogLevel:              opts.MySQLOptions.LogLevel,
		Logger:                nil,
	}

	return &Config{
		Log:             opts.Log,
		ServerRunConfig: opts.RPCServer,
		Register:        opts.ServiceRegister,
		Aggregation:     opts.AggregationOptions,
		Redis:           redisConfig,
		MySQL:           mysqlConfig,
		Availability:    opts.Availability,
		LoadBalance:     opts.LoadBalance,
		RPCClient:       opts.RPCClient,
	}, nil
}

// 创建服务通用配置
func (cfg *Config) CreateKitexOptions() (result []server.Option) {
	// 配置服务地址
	addrStr := net.JoinHostPort(cfg.ServerRunConfig.BindAddress, strconv.Itoa(cfg.ServerRunConfig.BindPort))
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		log.Fatal(err)
	}
	result = append(result, server.WithServiceAddr(addr))

	// 服务基本信息
	result = append(result, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: v1.ServiceName}))

	// 服务注册
	retryConfig := kretry.NewRetryConfig(
		kretry.WithMaxAttemptTimes(uint(cfg.Register.MaxAttemtTimes)),
		kretry.WithRetryDelay(cfg.Register.RetryDelay),
		kretry.WithObserveDelay(cfg.Register.ObserverDelay),
	)
	r, err := etcd.NewEtcdRegistryWithRetry(cfg.Register.Addrs, retryConfig)
	if err != nil {
		log.Fatal(err)
	}
	result = append(result, server.WithRegistry(r))

	// 限流
	result = append(result, server.WithLimit(&limit.Option{
		MaxConnections: cfg.ServerRunConfig.ConnectionLimit,
		MaxQPS:         cfg.ServerRunConfig.QPSLimit,
	}))

	// 设置优雅关停等待时间
	result = append(result, server.WithExitWaitTime(cfg.ServerRunConfig.ExitWaitTime))

	// 设置最大空闲时间
	result = append(result, server.WithMaxConnIdleTime(cfg.ServerRunConfig.MaxConnIdleTime))

	// 多路复用
	if cfg.ServerRunConfig.EnableMuxConnection {
		result = append(result, server.WithMuxTransport())
	}

	// 链路追踪
	result = append(result, server.WithSuite(tracing.NewServerSuite()))
	return result
}

// 创建rpc客户端通用配置
func (cfg *Config) CreateRemoteClientOptions() (result []client.Option) {
	// 长连接配置
	idleConfig := connpool.IdleConfig{
		MaxIdlePerAddress: cfg.RPCClient.MaxIdlePerAddress,
		MaxIdleGlobal:     cfg.RPCClient.MaxIdleGlobal,
		MaxIdleTimeout:    cfg.RPCClient.MaxIdleTimeout,
		MinIdlePerAddress: cfg.RPCClient.MinIdlePerAddress,
	}
	result = append(result, client.WithLongConnection(idleConfig))

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

func (cfg *Config) BuildEtcdConfig() *kv.EtcdConfig {
	return &kv.EtcdConfig{
		Endpoints: cfg.Register.Addrs,
		Username:  cfg.Register.Username,
		Password:  cfg.Register.Password,
	}
}

func (cfg *Config) BuildAggregation() {
	if cfg.Aggregation.Producer == "redis" {
		r := kv.NewKVStore("redis", nil).(kv.RedisStore)
		aggregation.NewAnalytics(cfg.Aggregation, producer.NewRedisProducer(r, cfg.Aggregation.StorageExpirationTime))
	} else if cfg.Aggregation.Producer == "rocketmq" {
		rmqProducer := producer.NewRocketMQProducer(cfg.Aggregation.RMQProducerOptions.Endpoints, cfg.Aggregation.RMQProducerOptions.AccessKey, cfg.Aggregation.RMQProducerOptions.AccessSecret)
		aggregation.NewAnalytics(cfg.Aggregation, rmqProducer)
	}
}
