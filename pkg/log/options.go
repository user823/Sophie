package log

import (
	"encoding/json"
	"fmt"
	"github.com/user823/Sophie/pkg/app"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

const (
	levelFlag             = "log.level"
	disableCallerFlag     = "log.disable_caller"
	disableStacktraceFlag = "log.disable_stacktrace"
	outputPathsFlag       = "log.output_paths"
	errPathsFlag          = "log.error_output_paths"
	developmentFlag       = "log.development"
	aggregationFlag       = "log.aggregation"
	nameFlag              = "log.name"
	disableLoggerFlag     = "log.disableLogger"

	stdout = "stdout"
	stderr = "stderr"
)

// 用于配置logger
// 配置可通过命令行参数进行修改
type Option struct {
	OutputPaths       []string `json:"output_paths" mapstructure:"output_paths" `
	ErrorOutputPaths  []string `json:"error_output_paths" mapstructure:"error_output_paths"`
	Level             string   `json:"level" mapstructure:"level"`
	DisableCaller     bool     `json:"disable_caller" mapstructure:"disable_caller"`
	DisableStacktrace bool     `json:"disable_stacktrace" mapstructure:"disable_stacktrace"`
	Development       bool     `json:"development" mapstructure:"development"`
	Aggregation       bool     `json:"aggregation" mapstructure:"aggregation"`
	DisableLogger     bool     `json:"disable_logger" mapstructure:"disable_logger"`
	Name              string   `json:"name" mapstructure:"name"`
}

func DefaultOptions() *Option {
	return &Option{
		Development:       false,
		OutputPaths:       []string{stdout},
		ErrorOutputPaths:  []string{stderr},
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Aggregation:       false,
		DisableLogger:     false,
	}
}

func (o *Option) Validate() []error {
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

func (o *Option) Flags() app.FlagGroup {
	fg := app.NewFlagGroup()
	fs := fg.FlagSet("logger")
	fs.StringVar(&o.Level, levelFlag, o.Level, "Mininum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, disableCallerFlag, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, disableStacktraceFlag,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringSliceVar(&o.OutputPaths, outputPathsFlag, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, errPathsFlag, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(&o.Aggregation, aggregationFlag, o.Aggregation, "Enable log aggregation, and store log records in Sophie-log.")
	fs.BoolVar(&o.DisableLogger, disableLoggerFlag, o.DisableLogger, "Disable logger in app.")
	fs.BoolVar(
		&o.Development,
		developmentFlag,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, nameFlag, o.Name, "The name of the logger.")
	return fg
}

func (o *Option) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
