package file

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/spf13/viper"
	v1 "github.com/user823/Sophie/api/domain/file/v1"
	"github.com/user823/Sophie/api/thrift/file/v1/fileservice"
	"github.com/user823/Sophie/internal/pkg/support"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/shutdown"
)

type FileServer struct {
	Gs              *shutdown.GracefulShutdown
	ShutdownTimeout time.Duration
	Analytics       *aggregation.Analytics
	Server          server.Server
	ServerConfig    *Config
}

func createGatewayServer(cfg *Config) (*FileServer, error) {
	// 创建基础组件
	gs := shutdown.NewGracefulShutdownInstance(v1.ServiceName)
	gs.SetErrHandler(&utils.GsLogErrHandler{})
	gs.AddShutdownManagers(shutdown.DefaultShutdownManager())
	gs.SetInOrder()

	if cfg.Log.Aggregation {
		cfg.BuildAggregation()
	}

	generalOpts := cfg.CreateKitexOptions()
	srv := fileservice.NewServer(new(FileServiceImpl), generalOpts...)

	return &FileServer{
		Gs:              gs,
		ShutdownTimeout: cfg.ServerRunConfig.ExitWaitTime,
		Analytics:       aggregation.GetAnalytics(),
		Server:          srv,
		ServerConfig:    cfg,
	}, nil
}

func (s *FileServer) PrepareRun() *FileServer {
	ctx, cancel := context.WithCancel(context.Background())

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

	// 最后关闭服务
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: rpc server is closing...", msg)

		if err := s.Server.Stop(); err != nil {
			log.Warnf("Shutdown rpc server error: %s", err.Error())
		}
		return nil
	})

	// 启动各种定时任务
	viperTask.Start()
	go kv.KeepConnection(ctx, s.ServerConfig.Redis)
	return s
}

func (s *FileServer) Run() error {
	// 启动各种服务组件
	if s.ServerConfig.Log.Aggregation {
		s.Analytics.Start()
	}

	// 开启优雅关停监听
	if err := s.Gs.Start(); err != nil {
		return err
	}

	tracingOpts := []provider.Option{
		provider.WithServiceName(v1.ServiceName),
		provider.WithEnableTracing(s.ServerConfig.Availability.TraceEnable),
		provider.WithExportEndpoint(s.ServerConfig.Availability.TraceEndpoint),
		provider.WithInsecure(),
		provider.WithEnableMetrics(s.ServerConfig.Availability.MetricEnable),
	}

	p := provider.NewOpenTelemetryProvider(tracingOpts...)
	defer p.Shutdown(context.Background())

	return s.Server.Run()
}
