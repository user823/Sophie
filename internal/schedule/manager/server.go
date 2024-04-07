package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/server"
	"github.com/google/uuid"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/spf13/viper"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	v12 "github.com/user823/Sophie/api/thrift/schedule/v1"
	"github.com/user823/Sophie/api/thrift/schedule/v1/jobservice"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/internal/pkg/support"
	"github.com/user823/Sophie/internal/schedule/locker"
	"github.com/user823/Sophie/internal/schedule/manager/loadbalance"
	"github.com/user823/Sophie/internal/schedule/manager/rpc"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/internal/schedule/store"
	"github.com/user823/Sophie/internal/schedule/store/mysql"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/log/aggregation"
	"github.com/user823/Sophie/pkg/shutdown"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// manager 节点
type ScheduleManager struct {
	ID              string
	Gs              *shutdown.GracefulShutdown
	ShutdownTimeout time.Duration
	Analytics       *aggregation.Analytics
	Server          server.Server
	ServerConfig    *Config
	etcdCli         kv.EtcdStore
	etcdNodePool    *models.EtcdNodePool
	lb              loadbalance.LoadBalancer
}

func createGatewayServer(cfg *Config) (*ScheduleManager, error) {
	// 创建基础组件
	gs := shutdown.NewGracefulShutdownInstance(v1.ServiceName)
	gs.SetErrHandler(&utils.GsLogErrHandler{})
	gs.AddShutdownManagers(shutdown.DefaultShutdownManager())
	gs.SetInOrder()

	if cfg.Log.Aggregation {
		cfg.BuildAggregation()
	}

	etcdStore, err := kv.NewEtcdClient(cfg.BuildEtcdConfig())
	if err != nil {
		log.Fatal("etcd client init failed")
	}
	lb := cfg.LoadBalance.CreateLoadBalancer()

	generalOpts := cfg.CreateKitexOptions()
	serviceImpl := new(JobServiceImpl)
	etcdNodePool := models.NewEtcdNodePool(etcdStore)
	serviceImpl.etcdNodePool = etcdNodePool
	serviceImpl.lb = lb
	srv := jobservice.NewServer(serviceImpl, generalOpts...)

	return &ScheduleManager{
		ID:              uuid.NewString(),
		Gs:              gs,
		ShutdownTimeout: cfg.ServerRunConfig.ExitWaitTime,
		Analytics:       aggregation.GetAnalytics(),
		Server:          srv,
		ServerConfig:    cfg,
		etcdCli:         etcdStore,
		etcdNodePool:    etcdNodePool,
		lb:              lb,
	}, nil
}

func (s *ScheduleManager) PrepareRun() *ScheduleManager {
	// 初始化服务组件
	if mycli, err := mysql.GetMySQLFactoryOr(s.ServerConfig.MySQL); err != nil {
		log.Fatal(err)
	} else {
		store.SetClient(mycli)
	}
	// 初始化死亡队列
	models.InitDeadJobQueue()
	// 初始化rpc
	rpcOpts := s.ServerConfig.CreateRemoteClientOptions()
	rpc.Init(rpcOpts)
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

	redisTask := support.GoTask{
		Run: func(ctx context.Context) (interface{}, error) {
			kv.KeepConnection(ctx, s.ServerConfig.Redis)
			return nil, nil
		},
	}
	tasks = append(tasks, &redisTask)

	// 节点状态监控
	watchNodeTask := support.GoTask{
		Run: s.WatchNodes,
	}
	tasks = append(tasks, &watchNodeTask)

	// 死亡队列处理
	deadJobsTask := support.GoTask{
		Run: s.DeadJobs,
	}
	tasks = append(tasks, &deadJobsTask)

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
	s.Gs.AddShutdownCallbacks(func(s string) error {
		log.Infof("%s: redis is closing...", s)
		redisTask.Shutdown(support.NOW)
		return nil
	})
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: closing worker nodes watch...", msg)
		watchNodeTask.Shutdown(support.NOW)
		return nil
	})

	// 死亡队列回调
	s.Gs.AddShutdownCallbacks(func(msg string) error {
		log.Infof("%s: closing dead jobs...", msg)
		deadJobsTask.Shutdown(support.NOW)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		result := deadJobsTask.Output(ctx)
		if !result.OK {
			return errors.New("get dead jobs failed")
		}

		// 将所有剩下的任务状态设置为暂停
		if jobs, ok := result.Data.([]any); ok {
			for i := range jobs {
				if job, ok := jobs[i].(*v1.SysJob); ok {
					job.Status = v1.PAUSE
					err := store.Client().Jobs().UpdateJob(context.Background(), job, &api.UpdateOptions{})
					if err != nil {
						log.Error("设置任务状态失败: %s", err.Error())
					}
				}
			}
		}
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

	// 启动各种定时任务
	for i := range tasks {
		tasks[i].Start()
	}
	return s
}

func (s *ScheduleManager) Run() error {
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

// node 节点监控
func (s *ScheduleManager) WatchNodes(ctx context.Context) (any, error) {
	cli := s.etcdCli.LowLevel().(*clientv3.Client)
	// 从最新版本开始监控
	var curRevision int64
	for {
		resp, err := cli.Get(context.Background(), models.NodesPrefix, clientv3.WithPrefix())
		if err != nil {
			continue
		}
		curRevision = resp.Header.Revision + 1
		break
	}

	ch := cli.Watch(ctx, models.NodesPrefix, clientv3.WithPrefix(), clientv3.WithRev(curRevision), clientv3.WithPrevKV())
	for watchResp := range ch {
		for _, event := range watchResp.Events {
			var nodeInfo models.Node

			switch event.Type {
			case clientv3.EventTypePut:
				if event.IsModify() {
					if err := json.Unmarshal(event.Kv.Value, &nodeInfo); err != nil {
						log.Error(err)
						continue
					}
					s.Seize(nodeInfo)
				}

			case clientv3.EventTypeDelete:
				if err := json.Unmarshal(event.PrevKv.Value, &nodeInfo); err != nil {
					log.Error(err)
					continue
				}
				nodeInfo.Status = v1.NODE_STATUS_OFF
				s.Seize(nodeInfo)
			}
		}
	}
	return nil, nil
}

func (s *ScheduleManager) Seize(node models.Node) {
	// 获取目标节点上运行的任务
	jobIds := s.etcdNodePool.JobsOnNode(node.Id)
	// 读取任务
	jobs, _, err := store.Client().Jobs().SelectJobByIds(context.Background(), jobIds, &api.GetOptions{})
	if err != nil {
		log.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 先抢锁, 抢到锁的manager负责调度
	if !s.scheduleTryLock(ctx, node.Id) {
		return
	}

	switch node.Status {
	case v1.NODE_STATUS_OFF:
		// 节点退出， 删除节点上的任务
		for i := range jobs {
			if err := s.etcdNodePool.Delete(jobs[i].JobId); err != nil {
				log.Error(err)
			}
			s.schedule(ctx, jobs[i])
		}
	case v1.NODE_STATUS_PAUSE:
		// 节点挂起，将目标节点上每个任务设置成暂停状态
		for i := range jobs {
			jobs[i].Status = v1.PAUSE
			if err := store.Client().Jobs().UpdateJob(ctx, jobs[i], &api.UpdateOptions{}); err != nil {
				log.Error(err)
			}
		}
	}
}

// 处理调度失败的任务
func (s *ScheduleManager) DeadJobs(ctx context.Context) (any, error) {
	dq := models.DeadJobQueue
	// 调度失败的任务处于低优先级，5秒调度一次
	delayTime := 5 * time.Second

	for {
		// 取出队首
		if !dq.Empty() {
			job, ok := dq.Poll().(*v1.SysJob)
			if !ok {
				continue
			}

			// 首先检查任务状态
			jobinfo, err := store.Client().Jobs().SelectJobById(ctx, job.JobId, &api.GetOptions{})
			// 查不到任务信息直接丢弃
			if err != nil || jobinfo == nil {
				continue
			}
			// 任务状态暂停直接丢弃
			if jobinfo.Status == v1.PAUSE {
				continue
			}
			// 任务已经被调度直接丢弃
			if _, ok = s.etcdNodePool.Exists(job.JobId); ok {
				continue
			}

			// 调度任务
			s.schedule(ctx, job)

			select {
			case <-ctx.Done():
				remain := dq.AllItems()
				return remain, nil
			case <-time.After(delayTime):
			}
		}
	}
}

// 抢占更新信息的锁
func (s *ScheduleManager) scheduleTryLock(ctx context.Context, nodeId string) bool {
	cli := s.etcdCli.LowLevel().(*clientv3.Client)
	res, err := cli.Grant(ctx, 5)
	if err != nil {
		log.Println(err)
		return false
	}
	resc, err := cli.KeepAlive(ctx, res.ID)
	if err != nil {
		log.Println(err)
		return false
	}

	go func(ctx context.Context) {
		for {
			select {
			case kresp := <-resc:
				if kresp != nil {
					log.Println("续租成功，LeaseID: ", kresp.ID)
				} else if resc == nil {
					log.Println("续租失败")
					return
				}
			case <-ctx.Done():
				if _, err := cli.Revoke(context.TODO(), res.ID); err != nil {
					log.Println(err)
				}
				return
			}
			time.Sleep(1 * time.Second)
		}
	}(ctx)

	key := fmt.Sprintf("%s:%s", locker.JOB_SCHEDULE, nodeId)
	txn := cli.Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, nodeId, clientv3.WithLease(res.ID))).
		Else(clientv3.OpGet(key))
	txnResp, err := txn.Commit()
	if err != nil {
		log.Println(err)
		return false
	}

	if !txnResp.Succeeded {
		return false
	} else {
		return true
	}
}

func (s *ScheduleManager) schedule(ctx context.Context, job *v1.SysJob) {
	strId := strconv.Itoa(int(job.JobId))
	picker := s.lb.GetPicker(s.etcdNodePool)

	failTimes := 0
	success := false

	// 每个任务尝试3次
	for failTimes < 3 {
		node := picker.Next(ctx, strId)
		// 当前节点不可用
		if node == nil {
			failTimes++
			continue
		}

		resp, err := rpc.Remoting.CreateJob(ctx, v12.SysJob2JobInfo(job), callopt.WithHostPort(node.Address().String()))
		// rpc调用失败
		if err != nil {
			failTimes++
			continue
		}
		success = true
		if resp.Code != code.SUCCESS {
			log.Errorf("死亡队列调度失败: %s", resp.Msg)
			break
		}
		log.Debugf("%d 任务调度成功", job.JobId)
	}

	// rpc 仍然调度失败，重新进入死亡队列末尾
	if !success {
		models.DeadJobQueue.Add(job.JobId, job)
	}
}
