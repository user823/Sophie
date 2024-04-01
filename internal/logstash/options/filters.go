package options

import flag "github.com/spf13/pflag"

type FilterOptions struct {
	Level            int8     `json:"log_level" mapstructure:"log_level"`
	LoggerWhiteList  []string `json:"logger_whitelist" mapstructure:"logger_whitelist"`
	LoggerBlackList  []string `json:"logger_blacklist" mapstructure:"logger_blacklist"`
	TimeFilterEnable bool     `json:"time_filter" mapstructure:"time_filter"`
}

func NewFilterOptions() *FilterOptions {
	return &FilterOptions{
		Level:            0,
		LoggerBlackList:  []string{},
		LoggerWhiteList:  []string{},
		TimeFilterEnable: false,
	}
}

func (o *FilterOptions) Validate() error { return nil }

func (o *FilterOptions) AddFlags(fs *flag.FlagSet) {
	fs.Int8Var(&o.Level, "log_level", o.Level, "set log level filter")
	fs.StringSliceVar(&o.LoggerWhiteList, "logger_whitelist", o.LoggerWhiteList, "set logger white list")
	fs.StringSliceVar(&o.LoggerBlackList, "logger_blacklist", o.LoggerBlackList, "set logger black list")
	fs.BoolVar(&o.TimeFilterEnable, "time_enable", o.TimeFilterEnable, "enable time filter")
}
