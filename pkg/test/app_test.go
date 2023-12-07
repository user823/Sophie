package test

import (
	"fmt"
	"github.com/user823/Sophie/pkg/app"
	"testing"
)

func TestApp(t *testing.T) {
	subCommand := app.NewCommand("testsub", "used to test subCommand",
		app.WithCommandRunFunc(func(args []string) error {
			fmt.Println("---开始打印参数---")
			fmt.Println(args)
			return nil
		}))

	myApp := app.NewApp("myApp",
		app.WithCommands(subCommand),
		app.WithVerbose(),
		app.WithConfigurable(),
		app.WithDefaultArgsValidation(),
		app.WithDescription("This is my first app"),
		app.WithRunFunc(run("myApp")))
	myApp.Run()
}

func TestSubCommand(t *testing.T) {
	subCommand := app.NewCommand("testsub", "used to test subCommand",
		app.WithCommandRunFunc(func(args []string) error {
			fmt.Println("---开始打印参数---")
			fmt.Println(args)
			return nil
		}))
	if err := subCommand.CobraCommand().Execute(); err != nil {
		fmt.Println(err.Error())
	}
}

func run(name string) app.RunFunction {
	return func(n string) error {
		fmt.Printf("--- %s run success ---\n", name)
		return nil
	}
}
