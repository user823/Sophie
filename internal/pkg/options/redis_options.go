package options

import (
	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/db/kv"
)

type RedisOptions struct {
	Addrs                 []string `json:"addrs"                    mapstructure:"addrs"`
	Username              string   `json:"username"                 mapstructure:"username"`
	Password              string   `json:"password"                 mapstructure:"password"`
	Database              int      `json:"database"                 mapstructure:"database"`
	MasterName            string   `json:"master_name"              mapstructure:"master_name"`
	MaxIdle               int      `json:"optimisation_max_idle"    mapstructure:"optimisation_max_idle"`
	MaxActive             int      `json:"optimisation_max_active"  mapstructure:"optimisation_max_active"`
	Timeout               int      `json:"timeout"                  mapstructure:"timeout"`
	EnableCluster         bool     `json:"enable_cluster"           mapstructure:"enable_cluster"`
	UseSSL                bool     `json:"use_ssl"                  mapstructure:"use_ssl"`
	SSLInsecureSkipVerify bool     `json:"ssl_insecure_skip_verify" mapstructure:"ssl_insecure_skip_verify"`
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addrs:                 []string{"127.0.0.1:6379"},
		Username:              "sophie",
		Password:              "123456",
		Database:              0,
		MasterName:            "",
		MaxIdle:               2000,
		MaxActive:             4000,
		Timeout:               0,
		EnableCluster:         false,
		UseSSL:                false,
		SSLInsecureSkipVerify: false,
	}
}

func (o *RedisOptions) Validate() error { return nil }

func (o *RedisOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringSliceVar(&o.Addrs, "redis.addrs", o.Addrs, "A set of redis address(format: 127.0.0.1:6379).")
	fs.StringVar(&o.Username, "redis.username", o.Username, "Username for access to redis service.")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Optional auth password for Redis db.")

	fs.IntVar(&o.Database, "redis.database", o.Database, ""+
		"By default, the database is 0. Setting the database is not supported with redis cluster. "+
		"As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.")

	fs.StringVar(&o.MasterName, "redis.master_name", o.MasterName, "The name of master redis instance.")

	fs.IntVar(&o.MaxIdle, "redis.optimisation_max_idle", o.MaxIdle, ""+
		"This setting will configure how many connections are maintained in the pool when idle (no traffic). "+
		"Set the --redis.optimisation-max-active to something large, we usually leave it at around 2000 for "+
		"HA deployments.")

	fs.IntVar(&o.MaxActive, "redis.optimisation_max_active", o.MaxActive, ""+
		"In order to not over commit connections to the Redis server, we may limit the total "+
		"number of active connections to Redis. We recommend for production use to set this to around 4000.")

	fs.IntVar(&o.Timeout, "redis.timeout", o.Timeout, "Timeout (in seconds) when connecting to redis service.")

	fs.BoolVar(&o.EnableCluster, "redis.enable_cluster", o.EnableCluster, ""+
		"If you are using Redis cluster, enable it here to enable the slots mode.")

	fs.BoolVar(&o.UseSSL, "redis.use_ssl", o.UseSSL, ""+
		"If set, IAM will assume the connection to Redis is encrypted. "+
		"(use with Redis providers that support in-transit encryption).")

	fs.BoolVar(&o.SSLInsecureSkipVerify, "redis.ssl_insecure_skip_verify", o.SSLInsecureSkipVerify, ""+
		"Allows usage of self-signed certificates when connecting to an encrypted Redis database.")
}

func (o *RedisOptions) BuildRdsConfig() *kv.RedisConfig {
	return &kv.RedisConfig{
		Addrs:                 o.Addrs,
		MasterName:            o.MasterName,
		Username:              o.Username,
		Password:              o.Password,
		Database:              o.Database,
		MaxIdle:               o.MaxIdle,
		MaxActive:             o.MaxActive,
		Timeout:               o.Timeout,
		EnableCluster:         o.EnableCluster,
		UseSSL:                o.UseSSL,
		SSLInsecureSkipVerify: o.SSLInsecureSkipVerify,
	}
}
