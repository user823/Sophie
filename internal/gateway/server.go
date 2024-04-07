package gateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	v1 "github.com/user823/Sophie/api/domain/gateway/v1"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/user823/Sophie/internal/gateway/router"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/pkg/support"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/shutdown"
)

type GateWayServer struct {
	Gs                           *shutdown.GracefulShutdown
	ShutdownTimeout              time.Duration
	Analytics                    *aggregation.Analytics
	InsecureServer, SecureServer *server.Hertz
	ServerConfig                 *Config
}

func createGatewayServer(config *Config) (*GateWayServer, error) {
	// 创建基础组件(redis, mysql ...)
	gs := shutdown.NewGracefulShutdownInstance(v1.ServiceName)
	gs.SetErrHandler(shutdown.EmptyErrHandler{})
	gs.AddShutdownManagers(shutdown.DefaultShutdownManager())
	gs.SetInOrder()

	if config.Log.Aggregation {
		config.BuildAggregation()
	}

	var insecureServer, secureServer *server.Hertz
	generalOpts := config.CreateHertzOptions()

	// gateway链路追踪
	if config.Availability.TraceEnable {
		tracer, cfg := hertztracing.NewServerTracer()
		insecureServer = server.New(
			append(generalOpts, tracer, server.WithHostPorts(config.InsecureServing.Addr))...,
		)
		secureServer = server.New(
			append(generalOpts, tracer, server.WithHostPorts(config.SecureServing.Addr),
				server.WithTLS(config.SecureServing.TLSConfig),
			)...,
		)
		insecureServer.Use(hertztracing.ServerMiddleware(cfg))
		secureServer.Use(hertztracing.ServerMiddleware(cfg))
	} else {
		insecureServer = server.New(append(generalOpts, server.WithHostPorts(config.InsecureServing.Addr))...)
		secureServer = server.New(append(generalOpts, server.WithHostPorts(config.SecureServing.Addr),
			server.WithTLS(config.SecureServing.TLSConfig),
		)...)
	}

	// 创建系统运行服务
	return &GateWayServer{
		Gs:              gs,
		ShutdownTimeout: config.ServerRunOptions.ExitWaitTime,
		Analytics:       aggregation.GetAnalytics(),
		InsecureServer:  insecureServer,
		SecureServer:    secureServer,
		ServerConfig:    config,
	}, nil
}

func (s *GateWayServer) PrepareRun() *GateWayServer {
	ctx, cancel := context.WithCancel(context.Background())

	// hertz环境准备
	s.HertzSetting()

	// 安装路由信息
	routerOpts := []router.Option{
		router.WithHealthz(s.ServerConfig.Healthz),
		router.WithMiddlewares(s.ServerConfig.Middlewares),
		router.WithJwtInfo(s.ServerConfig.Jwt),
		router.WithBaseAPI(s.ServerConfig.ServerRunOptions.BaseAPI),
		router.WithAddress(s.ServerConfig.InsecureServing.Addr),
	}
	router.InitRouter(s.InsecureServer, routerOpts...)
	router.InitRouter(s.SecureServer, routerOpts...)

	// 创建后台任务
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

	// 初始化服务组件
	rpcOpts := s.ServerConfig.CreateRemoteClientOptions()
	rpc.Init(rpcOpts)

	// 设置关停回调
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Info("Starting graceful shutdown...")
		cancel()
		return nil
	})
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: viper service is closing...", msg)
		viperTask.Shutdown(support.NOW)
		return nil
	})
	if s.ServerConfig.Log.Aggregation {
		s.Gs.AddShutdownCallbacks(func(msg string) error {
			log.Infof("%s: log aggregation is closing...", msg)
			aggregation.GetAnalytics().Stop()
			return nil
		})
	}

	// 最后关闭server
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: http/https server is closing...", msg)
		tCtx, tCancel := context.WithTimeout(ctx, s.ShutdownTimeout)
		defer tCancel()

		if err := s.SecureServer.Shutdown(tCtx); err != nil {
			log.Warnf("Shutdown secure server error: %s", err.Error())
		}
		if err := s.InsecureServer.Shutdown(tCtx); err != nil {
			log.Warnf("Shutdown insecure server error: %s", err.Error())
		}
		return nil
	})

	// 启动各种后台任务
	viperTask.Start()
	go kv.KeepConnection(ctx, s.ServerConfig.Redis)

	// 睡眠2s
	time.Sleep(2 * time.Second)
	return s
}

func (s *GateWayServer) Run() error {
	// 启动各种服务组件
	if s.ServerConfig.Log.Aggregation {
		s.Analytics.Start()
	}

	// 开启优雅关停监听
	if err := s.Gs.Start(); err != nil {
		return err
	}

	// 最后执行必要服务检查
	if !kv.Connected() {
		log.Fatalf("redis connect failed")
	}

	tracingOpts := []provider.Option{
		provider.WithEnableTracing(s.ServerConfig.Availability.TraceEnable),
		provider.WithExportEndpoint(s.ServerConfig.Availability.TraceEndpoint),
		provider.WithInsecure(),
		provider.WithEnableMetrics(s.ServerConfig.Availability.MetricEnable),
	}

	// 运行服务
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		opts := append(tracingOpts, provider.WithServiceName(v1.ServiceName+"/http"))
		p := provider.NewOpenTelemetryProvider(opts...)
		defer p.Shutdown(context.Background())
		s.InsecureServer.Spin()
	}()
	go func() {
		defer wg.Done()
		opts := append(tracingOpts, provider.WithServiceName(v1.ServiceName+"/https"))
		p := provider.NewOpenTelemetryProvider(opts...)
		defer p.Shutdown(context.Background())
		s.SecureServer.Spin()
	}()

	if s.ServerConfig.Healthz {
		time.Sleep(1 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := s.ping(ctx); err != nil {
			log.Errorf("Perform health check error: %s", err.Error())
		}
	}
	wg.Wait()

	log.Info("http server exit success")
	return nil
}

// 如果开启健康检查，服务启动后自动执行一次ping命令进行检查
func (s *GateWayServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/health", s.ServerConfig.InsecureServing.Addr)
	for {
		// 构造健康检查请求
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The service has been deployed successfully.")
			resp.Body.Close()
			return nil
		}

		log.Info("Waiting for the service, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}

func (s *GateWayServer) HertzSetting() {
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
	s.SecureServer.SetCustomSignalWaiter(waitSignal)
	s.InsecureServer.SetCustomSignalWaiter(waitSignal)
}
