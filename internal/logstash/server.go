package logstash

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/internal/logstash/exporters"
	"github.com/user823/Sophie/internal/logstash/options"
	"github.com/user823/Sophie/internal/pkg/middleware"
	"github.com/user823/Sophie/internal/pkg/support"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/log/aggregation/consumer"
	"github.com/user823/Sophie/pkg/shutdown"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// 默认最大10 个RecordConsumer
	maxConsumers = 10
)

type LogStash struct {
	Gs              *shutdown.GracefulShutdown
	ShutdownTimeout time.Duration
	InsecureServer  *server.Hertz
	options         *options.Options
	exporters       []exporters.RecordExporter
	consumers       []aggregation.RecordConsumer
	// 拉取record 间隔
	pumpInternal time.Duration
}

func createGatewayServer(opts *options.Options) (*LogStash, error) {
	// 创建基础组件(redis, mysql ...)
	gs := shutdown.NewGracefulShutdownInstance(ServiceName)
	gs.SetErrHandler(shutdown.EmptyErrHandler{})
	gs.AddShutdownManagers(shutdown.DefaultShutdownManager())
	gs.SetInOrder()

	generalOpts := opts.CreateHertzOptions()
	addr := net.JoinHostPort(opts.InsecureServing.BindAddress, strconv.Itoa(opts.InsecureServing.BindPort))
	insecureServer := server.New(append(generalOpts, server.WithHostPorts(addr))...)

	return &LogStash{
		Gs:              gs,
		ShutdownTimeout: opts.ServerRunOptions.ExitWaitTime,
		InsecureServer:  insecureServer,
		options:         opts,
		exporters:       make([]exporters.RecordExporter, len(opts.Exporters)),
		pumpInternal:    opts.Timeout,
	}, nil
}

func (s *LogStash) PrepareRun() *LogStash {
	// hertz 环境准备
	s.HertzSetting()
	s.initRouter()

	// 初始化exporters
	s.options.InitExporterConfigs()
	cnt := 0
	for _, name := range s.options.Exporters {
		exporter, err := exporters.NewExporter(name)
		if err != nil {
			log.Error(err)
			continue
		}
		s.exporters[cnt] = exporter
		cnt++
	}
	s.exporters = s.exporters[:cnt]

	// 创建后台任务
	var tasks []*support.GoTask
	viperTask := support.GoTask{}
	viperTask.Run = func(ctx context.Context) (interface{}, error) {
		// 每5s 拉取一次配置信息
		for {
			select {
			case <-ctx.Done():
				return nil, nil
			default:
				viperTask.WaitForRunning(5 * time.Second)
			}
			if err := viper.WatchRemoteConfig(); err != nil {
				log.Debug("Read remoting config err: %s", err.Error())
			}
		}
	}
	tasks = append(tasks, &viperTask)

	redisTask := support.GoTask{}
	// redis后台任务
	if s.options.SubRedisOptions != nil {
		rdsCfg := s.options.SubRedisOptions.BuildRdsConfig()
		redisTask.Run = func(ctx context.Context) (interface{}, error) {
			kv.KeepConnection(ctx, rdsCfg)
			return nil, nil
		}
		tasks = append(tasks, &redisTask)
	}

	// 设置关停回调
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: viper service is closing...", msg)
		viperTask.Shutdown(support.NOW)
		return nil
	})

	if s.options.SubRMQOptions != nil {
		s.Gs.AddShutdownCallbacks(func(s string) error {
			log.Infof("%s: redis is closing...", s)
			redisTask.Shutdown(support.NOW)
			return nil
		})
	}

	// 启动各种后台任务
	for i := range tasks {
		tasks[i].Start()
	}

	// 睡眠2s
	time.Sleep(2 * time.Second)
	return s
}

func (s *LogStash) Run() error {
	// 初始化消费者（它必须等上面的后台任务准备好）
	s.InitConsumers()

	// 添加剩余的后台任务
	var remainTasks []support.GoTask
	pumpTask := support.GoTask{Run: s.Pump}
	remainTasks = append(remainTasks, pumpTask)

	s.Gs.AddShutdownCallbacks(func(s string) error {
		log.Infof("%s: pump task is closing...", s)
		pumpTask.Shutdown(support.NOW)
		return nil
	})

	// 最后关闭服务
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: http server is closing...", msg)
		tCtx, tCancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer tCancel()

		if err := s.InsecureServer.Shutdown(tCtx); err != nil {
			log.Warnf("Shutdown insecure server error: %s", err.Error())
		}
		return nil
	})

	for i := range remainTasks {
		remainTasks[i].Start()
	}

	// 开启优雅关停监听
	if err := s.Gs.Start(); err != nil {
		return err
	}

	log.Infof("http server listening on address=[%s]:[%d]", s.options.InsecureServing.BindAddress, s.options.InsecureServing.BindPort)
	// 运行服务
	s.InsecureServer.Spin()
	log.Info("http server exit success")
	return nil
}

func (s *LogStash) HertzSetting() {
	// 关闭hertz 自带的log
	hlog.SetOutput(io.Discard)

	// 关闭hertz 自带的优雅关停
	waitSignal := func(errCh chan error) error {
		select {
		case err := <-errCh:
			// error occurs, exit immediately
			return err
		}
	}
	s.InsecureServer.SetCustomSignalWaiter(waitSignal)
}

func (s *LogStash) initRouter() {
	middlewares := middleware.Get("recovery", "requestid")
	middlewares = append(middlewares, middleware.Get(s.options.ServerRunOptions.Middlewares...)...)

	h := s.InsecureServer.Group(s.options.ServerRunOptions.BaseAPI)

	if s.options.ServerRunOptions.Healthz {
		h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
			core.OK(ctx, "ok", nil)
		})
	}
}

func (s *LogStash) InitConsumers() {
	s.consumers = make([]aggregation.RecordConsumer, maxConsumers)
	cnt := 0

	if s.options.SubRMQOptions != nil {
		endpoints := strings.Join(s.options.SubRMQOptions.Endpoints, ";")
		cs := consumer.NewRocketMQConsumer(endpoints, s.options.SubRMQOptions.AccessKey, s.options.SubRMQOptions.AccessSecret)
		s.consumers[cnt] = cs
		cnt++
	}

	if s.options.SubRedisOptions != nil {
		rdsStore := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsConsumer := consumer.NewRedisConsumer(rdsStore)
		s.consumers[cnt] = rdsConsumer
		cnt++
	}

	if cnt == 0 {
		log.Fatal("you have not provide available record aggregations source")
	}
	s.consumers = s.consumers[:cnt]

	for i := range s.consumers {
		s.consumers[i].Connect()
	}
}

func (s *LogStash) Pump(ctx context.Context) (any, error) {
	ticker := time.NewTicker(s.pumpInternal)
	defer ticker.Stop()

	log.Info("Now run loop to clean data with consumers")
	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			wg.Add(len(s.consumers))
			for i := range s.consumers {
				go s.pump(&wg, s.consumers[i])
			}
		case <-ctx.Done():
			log.Info("stop purge loop")
			return nil, nil
		}
	}
}

func (s *LogStash) pump(wg *sync.WaitGroup, consumer aggregation.RecordConsumer) {
	defer wg.Done()
	if !consumer.Connect() {
		return
	}

	records := consumer.GetAndDeleteSet(context.Background())
	if len(records) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.options.Timeout)
	defer cancel()

	if len(s.exporters) > 0 {
		var ewg sync.WaitGroup
		ewg.Add(len(s.exporters))
		for i := range s.exporters {
			go func(exp exporters.RecordExporter) {
				defer ewg.Done()
				err := exp.WriteData(ctx, records)
				if err != nil {
					log.Errorf("exporter write data error: %s", err.Error())
				}
			}(s.exporters[i])
		}
		ewg.Wait()
	} else {
		log.Warnf("no exporter has found")
	}
}
