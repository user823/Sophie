// 各个模块构建App 服务实例通用框架
package app

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/utils"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

var (
	progressMessage = "======>"
)

type App struct {
	// 实例名
	name        string
	description string
	verbose     bool
	// 设置App是否可配置
	// 如果可配置则配置使用顺序为：命令行 > 环境变量 > 配置文件 > 配置中心
	configurable bool
	options      CliOptions
	runE         RunFunction
	commands     []*Command
	args         cobra.PositionalArgs
	cmd          *cobra.Command
	output       io.Writer
}

type RunFunction func(name string) error

type Option func(*App)

// 为App实例添加描述信息
func WithDescription(description string) Option {
	return func(app *App) {
		app.description = description
	}
}

// 设置为静默启动
func WithVerbose() Option {
	return func(app *App) {
		app.verbose = true
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
				if len(arg) > 0 {
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

// 添加子命令
func WithCommands(cmds ...*Command) Option {
	return func(app *App) {
		app.commands = append(app.commands, cmds...)
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
		SilenceUsage:  a.verbose,
		SilenceErrors: a.verbose,
		Args:          a.args,
	}
	cmd.SetOut(a.output)
	cmd.SetErr(a.output)
	fs := cmd.Flags()
	ds.InitFlags(fs)

	// 添加子命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.CobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.name))
	}

	if a.runE != nil {
		cmd.RunE = a.runCommand()
	}

	flagGroup := ds.NewFlagGroup()
	if a.options != nil {
		fg := a.options.Flags()
		flagGroup.Merge(fg)
		// 将options 中的flagset都添加进来
		for _, f := range fg.FlagSets() {
			fs.AddFlagSet(f)
		}
	}

	// 尝试从配置文件或者配置中心拉取配置
	if a.configurable {
		flagGroup.AddGlobalFlags(addConfigFlag(a.name))
	}

	flagGroup.AddGlobalFlags(addHelpFlag(a.name))
	fs.AddFlagSet(flagGroup.FlagSet("global"))
	addCmdTemplate(cmd, flagGroup)
	a.cmd = cmd
}

func (a *App) runCommand() func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		printWorkingDir()
		// 打印flag信息
		if a.verbose {
			ds.PrintFlags(a.cmd.Flags())
		}

		if a.configurable {
			// 通过命令行更新viper中的配置
			if err := viper.BindPFlags(a.cmd.Flags()); err != nil {
				log.Printf("Viper flag bind error: %s", err.Error())
				return err
			}

			// 将viper的配置导出到options中
			// options 通过mapstructure 来设置键值, 注意options中的未导出字段会被decoder忽略掉
			if err := viper.Unmarshal(a.options); err != nil {
				log.Printf("Viper option unmarshal error: %s", err.Error())
				return err
			}
		}

		if a.verbose {
			log.Printf("%v Starting %s ...", progressMessage, a.name)
			if a.configurable && viper.ConfigFileUsed() != "" {
				log.Printf("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
			}
		}

		if a.options != nil {
			if err := a.applyOptionRules(); err != nil {
				return err
			}
		}

		if a.runE != nil {
			return a.runE(a.name)
		}
		return nil
	}
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", "Error:", err)
		os.Exit(1)
	}
}

// 应用选项规则
func (a *App) applyOptionRules() error {
	if c, ok := a.options.(CompletableOptions); ok {
		if err := c.Complete(); err != nil {
			return err
		}
	}

	if v, ok := a.options.(ValidatableOptions); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	if p, ok := a.options.(PrintableOptions); ok {
		log.Printf("%v Config: `%s`", progressMessage, p.String())
	}
	return nil
}

// 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Printf("%v WorkingDir: %s", progressMessage, wd)
}

// 为cobra 命令添加模版化信息
func addCmdTemplate(cmd *cobra.Command, fg *ds.FlagGroup) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := utils.GetTermInfo(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		// 打印子命令信息
		fmt.Fprintf(cmd.OutOrStderr(), subCmdText(cmd, cols))
		// 分组打印flag信息
		fg.PrintFlags(cmd.OutOrStderr(), cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		// 打印子命令信息
		fmt.Fprintf(cmd.OutOrStderr(), subCmdText(cmd, cols))
		// 分组打印flag信息
		fg.PrintFlags(cmd.OutOrStdout(), cols)
	})
}

func subCmdText(cmd *cobra.Command, maxWidth int) string {
	if len(cmd.Commands()) == 0 {
		return ""
	}

	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	tw.Write([]byte("\n\nCommands:\n\n"))
	for _, command := range cmd.Commands() {
		fmt.Fprintf(tw, "\t%s\t%s\n", command.Name(), ds.LimitWidth(command.Short, maxWidth-2-len(command.Name())))
		//buf.WriteString(command.Name() + " " + command.Short + "\n")
	}
	tw.Flush()
	return buf.String()
}
