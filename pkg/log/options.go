package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/user823/Sophie/pkg/ds"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable_caller"
	flagDisableStacktrace = "log.disable_stacktrace"
	flagOutputPaths       = "log.output_paths"
	flagErrPaths          = "log.error_output_paths"
	flagDevelopment       = "log.development"
	flagAggregation       = "log.aggregation"
	flagName              = "log.name"
	flagDisableLogger     = "log.disableLogger"

	stdout = "stdout"
	stderr = "stderr"
)

// 用于配置logger
// 配置可通过命令行参数进行修改
type Options struct {
	OutputPaths       []string `json:"output_paths" mapstructure:"output_paths" `
	ErrorOutputPaths  []string `json:"error_output_paths" mapstructure:"error_output_paths"`
	Level             string   `json:"level" mapstructure:"level"`
	DisableCaller     bool     `json:"disable_caller" mapstructure:"disable_caller"`
	DisableStacktrace bool     `json:"disable_stacktrace" mapstructure:"disable_stacktrace"`
	Development       bool     `json:"development" mapstructure:"development"`
	Aggregation       bool     `json:"aggregation" mapstructure:"aggregation"`
	DisableLogger     bool     `json:"disable_logger" mapstructure:"disable_logger"`
	SkipCaller        int      `json:"skip_caller" mapstructure:"skip_caller"`
	Name              string   `json:"name" mapstructure:"name"`
}

func DefaultOptions() *Options {
	return &Options{
		Development:       false,
		OutputPaths:       []string{stdout},
		ErrorOutputPaths:  []string{stderr},
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Aggregation:       false,
		DisableLogger:     false,
		SkipCaller:        3,
	}
}

func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	for _, path := range o.OutputPaths {
		if path != stdout && !isValidPath(path) {
			errs = append(errs, fmt.Errorf("Target output dir {%s} is invalid!", path))
			break
		}
	}

	for _, path := range o.ErrorOutputPaths {
		if path != stdout && !isValidPath(path) {
			errs = append(errs, fmt.Errorf("Target error-output dir {%s} is invalid!", path))
			break
		}
	}
	return errs
}

func isValidPath(path string) bool {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	// 所给的文件不存在
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (o *Options) Flags() *ds.FlagGroup {
	fg := ds.NewFlagGroup()
	fs := fg.FlagSet("logger")
	fs.StringVar(&o.Level, flagLevel, o.Level, "Mininum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(&o.Aggregation, flagAggregation, o.Aggregation, "Enable log aggregation, and store log records in Sophie-log.")
	fs.BoolVar(&o.DisableLogger, flagDisableLogger, o.DisableLogger, "Disable logger in app.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
	return fg
}

func (o *Options) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Mininum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(&o.Aggregation, flagAggregation, o.Aggregation, "Enable log aggregation, and store log records in Sophie-log.")
	fs.BoolVar(&o.DisableLogger, flagDisableLogger, o.DisableLogger, "Disable logger in app.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
