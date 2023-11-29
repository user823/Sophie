// 各个模块构建App 服务实例通用框架
package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/user823/Sophie/pkg/utils"
	"io"
	"os"
)

var (
	progressMessage = "======>"
)

type App struct {
	// 实例名
	name         string
	description  string
	silence      bool
	configurable bool
	options      CliOptions
	runE         RunFunction
	commands     []*Command
	args         cobra.PositionalArgs
	cmd          *cobra.Command
	output       io.Writer
}

type RunFunction func(cmd *cobra.Command, args []string) error

type Option func(*App)

// 为App实例添加描述信息
func WithDescription(description string) Option {
	return func(app *App) {
		app.description = description
	}
}

// 设置为静默启动
func WithSilence() Option {
	return func(app *App) {
		app.silence = true
	}
}

// 添加命令行选项参数
func WithOptions(options CliOptions) Option {
	return func(app *App) {
		app.options = options
	}
}

// 添加App实例运行函数
func WithRunFunc(run RunFunction) Option {
	return func(app *App) {
		app.runE = run
	}
}

// 添加App位置参数校验
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// 添加默认位置参数校验
func WithDefaultArgsValidation() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 1 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		}
	}
}

// 配置输出
func WithAppOutput(w io.Writer) Option {
	return func(app *App) {
		app.output = w
	}
}

// App是否支持配置
func WithConfigurable() Option {
	return func(app *App) {
		app.configurable = true
	}
}

// 通过NewApp(xxx).Run()运行
func NewApp(name string, options ...Option) *App {
	app := &App{name: name, output: os.Stdout}

	for _, opt := range options {
		opt(app)
	}

	app.buildCommand()
	return app
}

// 根据选项配置构建app.cmd
func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:           a.name + " [command] [flags]",
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  a.silence,
		SilenceErrors: a.silence,
		Args:          a.args,
	}
	cmd.SetOut(a.output)
	cmd.SetErr(a.output)
	fs := cmd.Flags()
	InitFlags(fs)

	for _, command := range a.commands {
		cmd.AddCommand(command.cobraCommand())
	}

	cmd.SetHelpCommand(helpCommand(a.name))
	if a.runE != nil {
		cmd.RunE = a.runE
	}

	flagGroup := NewFlagGroup()
	if a.options != nil {
		fg := a.options.Flags()
		flagGroup.Merge(fg)
		for _, f := range fg {
			fs.AddFlagSet(f)
		}
	}

	// 设置App启动的默认参数
	SetDefaultConfig()
	if a.configurable {
		addConfigFlag(a.name, flagGroup.FlagSet("global"))
	}
	flagGroup.AddGlobalFlags(cmd.Name())
	fs.AddFlagSet(flagGroup.FlagSet("global"))
	addCmdTemplate(cmd, flagGroup)
	a.cmd = cmd
}

// 为App添加子命令
func (a *App) AddCommand(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
}

// 为cobra 命令添加模版化信息
func addCmdTemplate(cmd *cobra.Command, fg FlagGroup) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := utils.GetTermInfo(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		PrintFlags(cmd.OutOrStderr(), fg, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		PrintFlags(cmd.OutOrStdout(), fg, cols)
	})
}
