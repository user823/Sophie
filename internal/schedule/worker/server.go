package worker

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/server"
	"github.com/google/uuid"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/spf13/viper"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/api/thrift/schedule/v1/workerservice"
	"github.com/user823/Sophie/internal/pkg/support"
	"github.com/user823/Sophie/internal/schedule/locker"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/internal/schedule/predefined"
	"github.com/user823/Sophie/internal/schedule/registry"
	"github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/internal/schedule/store/mysql"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/shutdown"
	utils2 "github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/intutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	defaultLoadWeight = 1000
)

type ScheduleWorker struct {
	// 实例id
	ID string
	// 实例状态
	Status          string
	Gs              *shutdown.GracefulShutdown
	ShutdownTimeout time.Duration
	Analytics       *aggregation.Analytics
	Server          server.Server
	ServerConfig    *Config
	registry        registry.Registry
	etcdCli         kv.EtcdStore
	jobMp           models.JobMap
	createdTime     time.Time
}

func (s *ScheduleWorker) NodeInfo() models.Node {
	addrStr := net.JoinHostPort(s.ServerConfig.ServerRunConfig.BindAddress, strconv.Itoa(s.ServerConfig.ServerRunConfig.BindPort))
	// 默认负载 为每个节点总工作数减去节点上存在的任务数
	weight := defaultLoadWeight - models.GetJobNums()
	mode := v1.WORKER
	description := fmt.Sprintf("%s node worker, created at %v", s.ID, utils2.Time2Str(s.createdTime))
	return models.Node{
		Id:          s.ID,
		Status:      s.Status,
		Network:     "tcp",
		IPAddress:   addrStr,
		LoadWeight:  weight,
		Mode:        mode,
		Description: description,
	}
}

func createGatewayServer(cfg *Config) (*ScheduleWorker, error) {
	// 创建基础组件
	gs := shutdown.NewGracefulShutdownInstance(v1.ServiceName + " worker")
	gs.SetErrHandler(&utils.GsLogErrHandler{})
	gs.AddShutdownManagers(shutdown.DefaultShutdownManager())
	gs.SetInOrder()

	if cfg.Log.Aggregation {
		cfg.BuildAggregation()
	}

	etcdStore, err := kv.NewEtcdClient(cfg.Register.BuildEtcdConfig())
	if err != nil {
		log.Fatal("etcd client init failed")
	}
	r, err := registry.NewEtcdRegistryWithClient(etcdStore, cfg.Register)
	if err != nil {
		log.Fatal(err)
	}
	etcdNodePool := models.NewEtcdNodePool(etcdStore)

	nodeId := uuid.NewString()
	generalOpts := cfg.CreateKitexOptions()
	workerImpl := new(WorkerServiceImpl)
	workerImpl.jobMp = etcdNodePool
	workerImpl.nodeid = nodeId
	srv := workerservice.NewServer(workerImpl, generalOpts...)

	return &ScheduleWorker{
		ID:              nodeId,
		Status:          v1.NODE_STATUS_ON,
		Gs:              gs,
		ShutdownTimeout: cfg.ServerRunConfig.ExitWaitTime,
		Analytics:       aggregation.GetAnalytics(),
		Server:          srv,
		ServerConfig:    cfg,
		registry:        r,
		etcdCli:         etcdStore,
		jobMp:           etcdNodePool,
	}, nil
}

func (s *ScheduleWorker) PrepareRun() *ScheduleWorker {
	// 初始化服务组件
	if mysqlCli, err := mysql.GetMySQLFactoryOr(s.ServerConfig.MySQL); err != nil {
		log.Fatal(err)
	} else {
		store.SetClient(mysqlCli)
	}
	// 初始化调度器
	models.InitSchedule()
	// 初始化调用目标
	predefined.TargetsInit()
	// 初始化分布式锁
	etcdCli, ok := s.etcdCli.LowLevel().(*clientv3.Client)
	if !ok {
		log.Fatal("etcd client init failed")
	}
	session, err := concurrency.NewSession(etcdCli)
	if err != nil {
		log.Fatal("init etcd session error: %s", err.Error())
	}
	locker.InitMutex(session)

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

	registryTask := support.GoTask{
		Run: func(ctx context.Context) (interface{}, error) {
			node := s.NodeInfo()
			s.registry.KeepRegister(node)
			return nil, nil
		},
	}
	tasks = append(tasks, &registryTask)

	redisTask := support.GoTask{
		Run: func(ctx context.Context) (interface{}, error) {
			kv.KeepConnection(ctx, s.ServerConfig.Redis)
			return nil, nil
		},
	}
	tasks = append(tasks, &redisTask)

	cleanTask := support.GoTask{
		Run: s.WatchStatus,
	}
	tasks = append(tasks, &cleanTask)

	// 设置关停回调
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

	// 关闭服务注册
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: deregister service...", msg)
		s.Off()
		node := s.NodeInfo()
		return s.registry.Deregister(node)
	})

	// 关闭清理任务
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: closing clean task...", msg)
		cleanTask.Shutdown(support.NOW)
		return nil
	})

	// 关闭redis
	s.Gs.AddShutdownCallbacks(func(s string) error {
		log.Infof("%s: redis is closing...", s)
		redisTask.Shutdown(support.NOW)
		return nil
	})

	// 最后关闭服务
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: rpc server is closing...", msg)

		if err := s.Server.Stop(); err != nil {
			log.Warnf("Shutdown rpc server error: %s", err.Error())
		}
		return nil
	})

	// 开启后台服务
	for i := range tasks {
		tasks[i].Start()
	}
	return s
}

func (s *ScheduleWorker) Run() error {
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

	s.createdTime = time.Now()
	return s.Server.Run()
}

func (s *ScheduleWorker) Pause() {
	s.Status = v1.NODE_STATUS_PAUSE
	models.PauseAll()
}

func (s *ScheduleWorker) Off() {
	s.Status = v1.NODE_STATUS_OFF
	models.CleanJob()
}

// 监控任务, 定期清除不属于本节点的任务
func (s *ScheduleWorker) WatchStatus(ctx context.Context) (any, error) {
	delayTime := 5 * time.Second
	for {
		jobs := s.jobMp.JobsOnNode(s.ID)
		allJob := models.GetAllJobs()
		for i := range allJob {
			if !intutil.ContainsAnyInt64(allJob[i], jobs...) {
				models.DeleteJob(allJob[i])
			}
		}

		select {
		case <-ctx.Done():
			return nil, nil
		case <-time.After(delayTime):
		}
	}
}
