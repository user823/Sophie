package logstash

import (
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/internal/logstash/options"
	"github.com/user823/Sophie/pkg/app"
	"github.com/user823/Sophie/pkg/log"
)

const (
	ServiceName = "Sophie Logstash"
	commandDesc = `Sophie Logstash is a part of sophie-log-aggregation system, she can gather log information from multi-source and export to multi-target.`
)

func NewApp(srvname string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(srvname,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithConfigurable(),
		app.WithDefaultArgsValidation(),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunction {
	return func(basename string) error {
		// 初始化日志组件
		l, err := log.New(opts.Log)
		if err != nil {
			return err
		}
		log.SetGlobal(l.WithValues(api.LOG_SERVICE, basename))

		server, err := createGatewayServer(opts)
		if err != nil {
			return err
		}
		return server.PrepareRun().Run()
	}
}
