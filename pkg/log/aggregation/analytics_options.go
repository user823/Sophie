package aggregation

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"time"
)

type AnalyticsOptions struct {
	PoolSize                int           `json:"pool_size" mapstructure:"pool_size"`
	RecordsBufferSize       uint64        `json:"records_buffer_size" mapstructure:"records_buffer_size"`
	FlushInterval           uint64        `json:"flush_interval" mapstructure:"flush_interval"`
	StorageExpirationTime   time.Duration `json:"storage_expiration_time" mapstructure:"storage_expiration_time"`
	Enable                  bool          `json:"enable" mapstructure:"enable"`
	EnableDetailedRecording bool          `json:"enable_detailed_recording" mapstructure:"enable_detailed_recording"`
}

// 默认日志传送选项
func NewAnalyticsOptions() *AnalyticsOptions {
	return &AnalyticsOptions{
		Enable:                  true,
		PoolSize:                50,
		RecordsBufferSize:       1000,
		FlushInterval:           200,
		EnableDetailedRecording: true,
		StorageExpirationTime:   time.Duration(24) * time.Hour,
	}
}

func (o *AnalyticsOptions) Validate() error {
	if o == nil {
		return nil
	}
	if o.Enable && (o.FlushInterval < 1 || o.FlushInterval > 1000) {
		return fmt.Errorf("log-record flush_interval %v must be between 1 and 1000", o.FlushInterval)
	}
	return nil
}

func (o *AnalyticsOptions) AddFlags(fs *flag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.Enable, "log-record.enable", o.Enable, ""+
		"This sets the iam-authz-server to record analytics data.")

	fs.IntVar(&o.PoolSize, "log-record.pool_size", o.PoolSize,
		"Specify number of pool workers.")

	fs.Uint64Var(&o.RecordsBufferSize, "log-record.records_buffer_size", o.RecordsBufferSize,
		"Specifies buffer size for pool workers (size of each pipeline operation).")

	fs.BoolVar(&o.EnableDetailedRecording, "log-record.enable_detailed_recording", o.EnableDetailedRecording,
		"Enable detailed analytics at the key level.")

	fs.DurationVar(&o.StorageExpirationTime, "log-record.storage_expiration_time", o.StorageExpirationTime, ""+
		"Set to a value larger than the Pump's purge_delay. "+
		"This allows the log-record data to exist long enough in Redis to be processed by sophie-log system.")
}
