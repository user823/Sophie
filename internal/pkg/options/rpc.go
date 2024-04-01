package options

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

type RPCClientOptions struct {
	// 长连接相关
	MinIdlePerAddress int           `json:"min_idle_per_address" mapstructure:"min_idle_per_address"`
	MaxIdlePerAddress int           `json:"max_idle_per_address" mapstructure:"max_idle_per_address"`
	MaxIdleGlobal     int           `json:"max_idle_global" mapstructure:"max_idle_global"`
	MaxIdleTimeout    time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`

	// 超时控制
	ConnTimeout time.Duration `json:"conn_timeout" mapstructure:"conn_timeout"`
	RPCTimeout  time.Duration `json:"rpc_timeout" mapstructure:"rpc_timeout"`

	// 重试控制（异常重试、复用熔断器）
	MaxRetryTimes int    `json:"max_retry_times" mapstructure:"max_retry_times"`
	MaxDurationMS uint32 `json:"max_duration_ms" mapstructure:"max_duration_ms"`

	// 熔断率
	Circuitbreak float64 `json:"circuitbreak" mapstructure:"circuitbreak"`
	Minsample    int64   `json:"minsample" mapstructure:"minsample"`
}

type RPCServerOptions struct {
	// ip
	BindAddress string `json:"bind_address" mapstructure:"bind_address"`
	// port
	BindPort int `json:"bind_port"    mapstructure:"bind_port"`

	// 多路复用
	EnableMuxConnection bool `json:"enable_mux_connection" mapstructure:"enable_mux_connection"`

	// 连接闲置控制
	MaxConnIdleTime time.Duration `json:"max_conn_idle_time" mapstructure:"max_conn_idle_time"`

	// 限流相关
	QPSLimit        int `json:"qps_limit" mapstructure:"qps_limit"`
	ConnectionLimit int `json:"connection_limit" mapstructure:"connection_limit"`

	// 关停时间
	ExitWaitTime time.Duration `json:"exit_wait_time" mapstructure:"exit_wait_time"`
}

func NewRPCClientOptions() *RPCClientOptions {
	return &RPCClientOptions{
		MaxIdlePerAddress: 10,
		MaxIdleGlobal:     100,
		MaxIdleTimeout:    time.Minute,
		MinIdlePerAddress: 2,
		ConnTimeout:       3 * time.Second,
		RPCTimeout:        3 * time.Second,
		MaxRetryTimes:     2,
		MaxDurationMS:     0,
		Circuitbreak:      0.8,
		Minsample:         200,
	}
}

func (o *RPCClientOptions) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.MaxIdlePerAddress, "rpc_client.maxIdlePerAddress", o.MaxIdlePerAddress, ""+
		"Remoting invoke idle connections on each server address limit")
	fs.IntVar(&o.MaxIdleGlobal, "rpc_client.maxIdleGlobal", o.MaxIdleGlobal, ""+
		"Remoting invoke total idle connections limit")
	fs.DurationVar(&o.MaxIdleTimeout, "rpc_client.maxIdleTimeout", o.MaxIdleTimeout, ""+
		"Remoting invoke idle connection timeout")
	fs.IntVar(&o.MinIdlePerAddress, "rpc_client.minIdlePerAddress", o.MinIdlePerAddress, ""+
		"Remoting invoke min keep idle connections on each server address")
	fs.DurationVar(&o.ConnTimeout, "rpc_client.ConnTimeout", o.ConnTimeout, ""+
		"Remoting invoke connection timeout")
	fs.DurationVar(&o.RPCTimeout, "rpc_client.rpcTimeout", o.RPCTimeout, ""+
		"Remoting invoke timeout")
	fs.IntVar(&o.MaxRetryTimes, "rpc_client.maxRetryTimes", o.MaxRetryTimes, ""+
		"Remoting invoke max retry times (exclude first try)")
	fs.Uint32Var(&o.MaxDurationMS, "rpc_client.maxDurationMS", o.MaxDurationMS, ""+
		"Remoting invoke max retry accumulated time")
	fs.Float64Var(&o.Circuitbreak, "rpc_client.circuitbreak", o.Circuitbreak, ""+
		"Remoting invoke circuitbreak percent")
	fs.Int64Var(&o.Minsample, "rpc_client.minsample", o.Minsample, ""+
		"Remoting invoke circuitbreak recover sample num")
}

func NewRPCServerOptions() *RPCServerOptions {
	return &RPCServerOptions{
		EnableMuxConnection: true,
		MaxConnIdleTime:     30 * time.Minute,
		QPSLimit:            100,
		ConnectionLimit:     1000,
	}
}

func (o *RPCServerOptions) Validate() error {
	// rpc server option创建时不提供ip 和 addr的默认运行参数，防止多个微服务之间冲突
	// 需要从配置中心或者配置文件中读取
	if !utils.IsValidIP(o.BindAddress) {
		return fmt.Errorf("Error rpc server bind address %s, please use ipv4 or ipv6 ", o.BindAddress)
	}
	if o.BindPort < 1 || o.BindPort > 65535 {
		return fmt.Errorf("Error rpc server bind port %d must be between 1 and 65535", o.BindPort)
	}
	return nil
}

func (o *RPCServerOptions) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&o.EnableMuxConnection, "rpc_server.enable_muxconnection", o.EnableMuxConnection, ""+
		"Remoting server use mux tcp connection")
	fs.DurationVar(&o.MaxConnIdleTime, "rpc_server.maxconn_idletime", o.MaxConnIdleTime, ""+
		"Remoting server max idle connection timeout")
	fs.IntVar(&o.QPSLimit, "rpc_server.qps_limit", o.QPSLimit, ""+
		"Remoting server max QPS")
	fs.IntVar(&o.ConnectionLimit, "rpc_server.connection_timeout", o.ConnectionLimit, ""+
		"Remoting server connections limit")
	fs.StringVar(&o.BindAddress, "rpc_server.bind_ipaddress", o.BindAddress, ""+
		"Remoting server ip address")
	fs.IntVar(&o.BindPort, "rpc_server.bind_port", o.BindPort, ""+
		"Remoting server listening port")
}
