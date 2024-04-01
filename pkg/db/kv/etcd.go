package kv

import (
	"context"
	"crypto/tls"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type EtcdConfig struct {
	Endpoints             []string
	Username              string
	Password              string
	Timeout               int
	UseSSL                bool
	SSLInsecureSkipVerify bool
	// 发送、接受最大消息大小（Byte为单位）
	MaxSize int
	Logger  *zap.Logger
}

type EtcdClient struct {
	*clientv3.Client
	closed   atomic.Bool
	leaseIds sync.Map
}

func NewEtcdClient(config any) (EtcdStore, error) {
	etcdCli, err := connectEtcd(config)
	if err != nil {
		return nil, err
	}
	return &EtcdClient{Client: etcdCli}, nil
}

func connectEtcd(config any) (*clientv3.Client, error) {
	cfg, ok := config.(*EtcdConfig)
	if !ok {
		return nil, ErrConfigTypeInvalid
	}

	timeout := 3 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}

	// 连接超时
	dialTimeout := timeout

	// ping超时
	dialKeepAliveTimeout := 5 * timeout

	// tls config
	var tlsConfig *tls.Config
	if cfg.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: cfg.SSLInsecureSkipVerify,
		}
	}

	return clientv3.New(clientv3.Config{
		Endpoints:            cfg.Endpoints,
		Username:             cfg.Username,
		Password:             cfg.Password,
		DialKeepAliveTimeout: dialKeepAliveTimeout,
		DialTimeout:          dialTimeout,
		MaxCallSendMsgSize:   cfg.MaxSize,
		MaxCallRecvMsgSize:   cfg.MaxSize,
		TLS:                  tlsConfig,
		Logger:               cfg.Logger,
	})
}

// 修改客户端连接
func (e *EtcdClient) Connect(ctx context.Context, config any) {
	c, err := connectEtcd(config)
	if err != nil {
		log.Errorf("Connect to etcd error: %s", err.Error())
		return
	}
	e.Client = c
}

func (e *EtcdClient) Connected() bool {
	return !e.closed.Load()
}

func (e *EtcdClient) Disconnect() error {
	e.closed.Store(true)
	return e.Close()
}

func (e *EtcdClient) GetKey(ctx context.Context, key string) (string, error) {
	resp, err := e.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		return utils.B2s(resp.Kvs[0].Value), nil
	}
	return "", nil
}

func (e *EtcdClient) GetMultiKey(ctx context.Context, keys []string) ([]string, error) {
	result := make([]string, len(keys))
	found := false
	for i, key := range keys {
		resp, err := e.Get(ctx, key)
		if err != nil {
			result[i] = ""
			continue
		}
		if resp.Kvs[0].Value != nil {
			found = true
			result[i] = utils.B2s(resp.Kvs[0].Value)
		}
	}

	if !found {
		return []string{}, ErrKeyNotFound
	}
	return result, nil
}

func (e *EtcdClient) SetKey(ctx context.Context, key string, val string, expire int64) error {
	if expire > 0 {
		leaseResp, err := e.Grant(ctx, utils.NanoToSecond(expire))
		if err != nil {
			return err
		}
		e.leaseIds.Store(key, leaseResp.ID)
		_, err = e.Put(ctx, key, val, clientv3.WithLease(leaseResp.ID))
		return err
	}

	_, err := e.Put(ctx, key, val)
	return err
}

func (e *EtcdClient) SetExp(ctx context.Context, key string, expire int64) error {
	val, err := e.GetKey(ctx, key)
	if err != nil {
		return ErrKeyNotFound
	}
	return e.SetKey(ctx, key, val, expire)
}

func (e *EtcdClient) GetExp(ctx context.Context, key string) (int64, error) {
	leaseid, ok := e.leaseIds.Load(key)
	if !ok {
		return 0, ErrKeyNotFound
	}
	resp, err := e.TimeToLive(ctx, leaseid.(clientv3.LeaseID))
	if err != nil {
		return 0, err
	}
	return resp.TTL, nil
}

func (e *EtcdClient) GetKeys(ctx context.Context, prefix string) []string {
	resp, err := e.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return []string{}
	}
	result := make([]string, len(resp.Kvs))
	for i := range resp.Kvs {
		result[i] = utils.B2s(resp.Kvs[i].Key)
	}
	return result
}

func (e *EtcdClient) DeleteKey(ctx context.Context, key string) bool {
	_, err := e.Delete(ctx, key)
	if err != nil {
		log.Errorf("delete key error: %s", err.Error())
	}
	return err != nil
}

func (e *EtcdClient) DeleteAllKeys(ctx context.Context) bool {
	_, err := e.Delete(ctx, "", clientv3.WithPrefix())
	log.Errorf("delete key error: %s", err.Error())
	return err != nil
}

func (e *EtcdClient) GetKeysAndValues(ctx context.Context) map[string]string {
	resp, err := e.Get(ctx, "", clientv3.WithPrefix())
	result := make(map[string]string)
	if err != nil {
		return result
	}
	for i := range resp.Kvs {
		k, v := resp.Kvs[i].Key, resp.Kvs[i].Value
		result[utils.B2s(k)] = utils.B2s(v)
	}
	return result
}

func (e *EtcdClient) GetKeysAndValuesWithFilter(ctx context.Context, prefix string) map[string]string {
	resp, err := e.Get(ctx, prefix, clientv3.WithPrefix())
	result := make(map[string]string)
	if err != nil {
		return result
	}
	for i := range resp.Kvs {
		k, v := resp.Kvs[i].Key, resp.Kvs[i].Value
		result[utils.B2s(k)] = utils.B2s(v)
	}
	return result
}

func (e *EtcdClient) DeleteKeys(ctx context.Context, keys []string) bool {
	for i := range keys {
		_, err := e.Delete(ctx, keys[i])
		if err != nil {
			log.Error(err)
		}
	}
	return true
}

func (e *EtcdClient) GrantLease(ctx context.Context, ttl int64) (clientv3.LeaseID, error) {
	resp, err := e.Grant(ctx, ttl)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func (e *EtcdClient) PutWithLease(ctx context.Context, key, val string, id clientv3.LeaseID) error {
	_, err := e.Put(ctx, key, val, clientv3.WithLease(id))
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdClient) KeepLease(ctx context.Context, leaseId clientv3.LeaseID) error {
	ch, err := e.KeepAlive(ctx, leaseId)
	if err != nil {
		return err
	}

	go func() {
		log.Infof("start keepalive lease %x", leaseId)
		for range ch {
			select {
			case <-ctx.Done():
				log.Infof("stop keepalive lease: %s", leaseId)
				if _, err = e.Revoke(context.Background(), leaseId); err != nil {
					log.Error(err)
				}
				return
			case resp := <-ch:
				if resp == nil {
					log.Errorf("keep lease %d failed", leaseId)
					return
				}
			}
		}
	}()
	return nil
}

func (e *EtcdClient) LowLevel() any {
	return e.Client
}
