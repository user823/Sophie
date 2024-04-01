package worker

import (
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/app"
	"github.com/user823/Sophie/pkg/log"
)

const (
	commandDesc = `Sophie Schedule worker is used to run jobs on worker node.`
)

func NewApp(srvname string) *app.App {
	opts := NewOptions()
	application := app.NewApp(srvname,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithVerbose(),
		app.WithConfigurable(),
		app.WithDefaultArgsValidation(),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *Options) app.RunFunction {
	return func(basename string) error {
		cfg, err := CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		// 初始化日志组件
		l, err := log.New(cfg.Log)
		if err != nil {
			return err
		}
		log.SetGlobal(l.WithValues(api.LOG_SERVICE, basename))

		server, err := createGatewayServer(cfg)
		if err != nil {
			return err
		}
		return server.PrepareRun().Run()
	}
}
