package gen

import (
	"net"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/db/sql"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/log/aggregation/producer"
)

// 运行、创建服务必要配置
type Config struct {
	Log             *log.Options
	ServerRunConfig *options.RPCServerOptions
	Register        *RegisterInfo
	Aggregation     *aggregation.AnalyticsOptions
	Redis           *kv.RedisConfig
	MySQL           *sql.MysqlConfig
	Availability    *options.AvailabilityOptions
	Gen             *GenOptions
}

func CreateConfigFromOptions(opts *Options) (*Config, error) {
	// 服务发现
	register := &RegisterInfo{
		Addrs:           opts.ServiceRegister.Addrs,
		Username:        opts.ServiceRegister.Username,
		Password:        opts.ServiceRegister.Password,
		MaxAttemptTimes: opts.ServiceRegister.MaxAttemtTimes,
		ObserverDelay:   opts.ServiceRegister.ObserverDelay,
		RetryDelay:      opts.ServiceRegister.RetryDelay,
	}

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
		Register:        register,
		Aggregation:     opts.AggregationOptions,
		Redis:           redisConfig,
		MySQL:           mysqlConfig,
		Availability:    opts.Availability,
		Gen:             opts.GenOptions,
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

	// 服务发现
	retryConfig := retry.NewRetryConfig(
		retry.WithMaxAttemptTimes(uint(cfg.Register.MaxAttemptTimes)),
		retry.WithRetryDelay(cfg.Register.RetryDelay),
		retry.WithObserveDelay(cfg.Register.ObserverDelay),
	)
	r, err := etcd.NewEtcdRegistryWithRetry(cfg.Register.Addrs, retryConfig)
	if err != nil {
		log.Fatal(err)
	}
	result = append(result, server.WithRegistry(r))

	// 链路追踪
	result = append(result, server.WithSuite(tracing.NewServerSuite()))
	return result
}

type RegisterInfo struct {
	Addrs           []string
	Username        string
	Password        string
	MaxAttemptTimes int
	// 服务发现正常时延
	ObserverDelay time.Duration
	// 服务发现重试时延
	RetryDelay time.Duration
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
