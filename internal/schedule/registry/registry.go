package registry

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Registry interface {
	Register(node models.Node) error
	KeepRegister(node models.Node)
	Deregister(node models.Node) error
}

type etcdRegistry struct {
	client   kv.EtcdStore
	leaseTTL int64
	stop     chan struct{}
	lease    *leaseInfo
	retry    *retryConfig
}

type leaseInfo struct {
	leaseId clientv3.LeaseID
	ctx     context.Context
	cancel  context.CancelFunc
}

type retryConfig struct {
	MaxAttemtTimes int
	ObserverDelay  time.Duration
	RetryDelay     time.Duration
}

func NewEtcdRegistry(opt *options.ServiceDiscoverOptions) (Registry, error) {
	etcdConfig := &kv.EtcdConfig{
		Endpoints: opt.Addrs,
		Username:  opt.Username,
		Password:  opt.Password,
	}

	retryCfg := &retryConfig{
		MaxAttemtTimes: opt.MaxAttemtTimes,
		ObserverDelay:  opt.ObserverDelay,
		RetryDelay:     opt.RetryDelay,
	}

	etcdCli, err := kv.NewEtcdClient(etcdConfig)
	if err != nil {
		return nil, err
	}

	return &etcdRegistry{
		client:   etcdCli,
		leaseTTL: int64(opt.TTL),
		retry:    retryCfg,
		stop:     make(chan struct{}, 1),
	}, nil
}

func NewEtcdRegistryWithClient(client kv.EtcdStore, opt *options.ServiceDiscoverOptions) (Registry, error) {
	retryCfg := &retryConfig{
		MaxAttemtTimes: opt.MaxAttemtTimes,
		ObserverDelay:  opt.ObserverDelay,
		RetryDelay:     opt.RetryDelay,
	}
	return &etcdRegistry{
		client:   client,
		leaseTTL: int64(opt.TTL),
		retry:    retryCfg,
		stop:     make(chan struct{}, 1),
	}, nil
}

func (e *etcdRegistry) Register(node models.Node) error {
	leaseId, err := e.client.GrantLease(context.Background(), e.leaseTTL)
	if err != nil {
		return err
	}
	if err = e.register(node, leaseId); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.lease = &leaseInfo{
		leaseId: leaseId,
		ctx:     ctx,
		cancel:  cancel,
	}
	if err = e.keepalive(); err != nil {
		return err
	}
	return nil
}

// 执行一次注册
func (e *etcdRegistry) register(node models.Node, leaseId clientv3.LeaseID) error {
	val, err := jsoniter.Marshal(node)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return e.client.PutWithLease(ctx, models.ServiceKey(node.CacheKey()), utils.B2s(val), leaseId)
}

func (e *etcdRegistry) keepalive() error {
	return e.client.KeepLease(e.lease.ctx, e.lease.leaseId)
}

// 保持注册
func (e *etcdRegistry) KeepRegister(node models.Node) {
	var failedTimes int
	delay := e.retry.ObserverDelay
	// maxAttemtTimes 为0时无限重试
	for e.retry.MaxAttemtTimes == 0 || failedTimes < e.retry.MaxAttemtTimes {
		select {
		case _, ok := <-e.stop:
			if !ok {
				close(e.stop)
			}
			log.Infof("stop keep register node: %s", node.Id)
			return
		case <-time.After(delay):
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		val, err := e.client.GetKey(ctx, models.ServiceKey(node.CacheKey()))
		if err != nil {
			log.Warnf("keep register get %s failed with err: %v", models.ServiceKey(node.CacheKey()), err)
			delay = e.retry.RetryDelay
			failedTimes++
			continue
		}

		// lease 可能过期了，重新注册
		if val == "" {
			log.Infof("keep register: %s", models.ServiceKey(node.CacheKey()))
			delay = e.retry.RetryDelay
			err = e.Register(node)
			if err != nil {
				log.Warnf("keep register get %s failed with err: %v", models.ServiceKey(node.CacheKey()), err)
				failedTimes++
				continue
			}
			delay = e.retry.ObserverDelay
		}

		// 注册成功清空失败次数
		failedTimes = 0
	}
	log.Errorf("keep register service %s failed times:%d", models.ServiceKey(node.CacheKey()), failedTimes)
}

func (e *etcdRegistry) Deregister(node models.Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if !e.client.DeleteKey(ctx, models.ServiceKey(node.CacheKey())) {
		return errors.Errorf("etcd client delete key %s failed", models.ServiceKey(node.CacheKey()))
	}
	e.stop <- struct{}{}
	return nil
}
