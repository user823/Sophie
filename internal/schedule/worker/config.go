package worker

import (
	"net"
	"strconv"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
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
	Register        *options.ServiceDiscoverOptions
	Aggregation     *aggregation.AnalyticsOptions
	MySQL           *sql.MysqlConfig
	Redis           *kv.RedisConfig
	Availability    *options.AvailabilityOptions
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

	// mysql
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
		MySQL:           mysqlConfig,
		Availability:    opts.Availability,
		Redis:           redisConfig,
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

	// 链路追踪
	result = append(result, server.WithSuite(tracing.NewServerSuite()))
	return result
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
