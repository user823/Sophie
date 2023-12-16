package app

import (
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"strings"
)

const (
	helpFlag  = "help"
	helpShort = "h"
)

// 构建cobra帮助命令
func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command.",
		Long: `Help provides help for any command in the application.
Simply type ` + name + ` help [path to command] for full details.`,

		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				_ = c.Root().Usage()
			} else {
				cmd.InitDefaultHelpFlag() // make possible 'help' flag to be shown
				_ = cmd.Help()
			}
		},
	}
}

func addHelpFlag(name string) *flag.Flag {
	flag.BoolP(helpFlag, helpShort, false, fmt.Sprintf("Help for %s.", name))
	return flag.Lookup(helpFlag)
}

func addHelpCommandFlag(usage string, fs *flag.FlagSet) {
	fs.BoolP(
		helpFlag,
		helpShort,
		false,
		fmt.Sprintf("Help for the %s command.", strings.Split(usage, " ")[0]),
	)
}
