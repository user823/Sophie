package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// 用于构建子命令树
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	outPut   io.Writer
	commands []*Command
	runE     RunCommandFunc
}

type CommandOption func(command *Command)

type RunCommandFunc func([]string) error

func WithCommandOptions(options CliOptions) CommandOption {
	return func(command *Command) {
		command.options = options
	}
}

func WithCommandRunFunc(runE RunCommandFunc) CommandOption {
	return func(command *Command) {
		command.runE = runE
	}
}

func WithCommandOutput(w io.Writer) CommandOption {
	return func(command *Command) {
		command.outPut = w
	}
}

func NewCommand(usage, desc string, options ...CommandOption) *Command {
	cmd := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, opt := range options {
		opt(cmd)
	}
	return cmd
}

func (c *Command) AddCommand(cmd ...*Command) {
	c.commands = append(c.commands, cmd...)
}

func (c *Command) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}

	cmd.Flags().SortFlags = false
	for _, subcommand := range c.commands {
		cmd.AddCommand(subcommand.CobraCommand())
	}

	if c.runE != nil {
		cmd.RunE = c.runCommand()
	}

	if c.options != nil {
		for _, fs := range c.options.Flags() {
			cmd.Flags().AddFlagSet(fs)
		}
	}

	// 添加帮助命令flag
	addHelpCommandFlag(c.usage, cmd.Flags())
	return cmd
}

func (c *Command) runCommand() func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		if c.runE != nil {
			if err := c.runE(args); err != nil {
				fmt.Printf("%v %v\n", "Error: ", err)
				os.Exit(1)
			}
		}
		return nil
	}
}
