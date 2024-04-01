package registry

import (
	"context"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	jsoniter "github.com/json-iterator/go"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/pkg/options"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/utils"
	"time"
)

const (
	// 数据同步超时时间(s)

	EtcdResolver = "etcd-resolver"
)

// 实现了kitex discovery 接口
type etcdResolver struct {
	client kv.EtcdStore
}

func NewEtcdResolver(opt *options.ServiceDiscoverOptions) (discovery.Resolver, error) {
	etcdConfig := &kv.EtcdConfig{
		Endpoints: opt.Addrs,
		Username:  opt.Username,
		Password:  opt.Password,
	}

	etcdCli, err := kv.NewEtcdClient(etcdConfig)
	if err != nil {
		return nil, err
	}

	return &etcdResolver{etcdCli}, nil
}

func NewEtcdResolverWithClient(client kv.EtcdStore) (discovery.Resolver, error) {
	return &etcdResolver{client}, nil
}

func (e *etcdResolver) sync(prefix string) (result []models.Node) {
	ctx, cancel := context.WithTimeout(context.Background(), models.SyncTimeout*time.Second)
	defer cancel()
	res := e.client.GetKeysAndValuesWithFilter(ctx, prefix)
	result = make([]models.Node, 0, len(res))
	cnt := 0
	for _, v := range res {
		var node models.Node
		if err := jsoniter.Unmarshal(utils.S2b(v), &node); err != nil {
			continue
		}
		result = append(result, node)
		cnt++
	}
	result = result[:cnt]
	return result
}

// 获取所有在线worker节点
func (e *etcdResolver) onlineNodes(prefix string) []models.Node {
	res := e.sync(prefix)
	cnt := 0
	for i := range res {
		if res[i].Status == v1.NODE_STATUS_ON && res[i].Mode == v1.WORKER {
			res[cnt] = res[i]
			cnt++
		}
	}
	res = res[:cnt]
	return res
}

func (e *etcdResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	return target.ServiceName()
}

func (e *etcdResolver) Resolve(ctx context.Context, desc string) (discovery.Result, error) {
	nodes := e.onlineNodes(desc)
	if len(nodes) == 0 {
		return discovery.Result{}, errors.Errorf("no instance remains for %s", desc)
	}
	result := make([]discovery.Instance, 0, len(nodes))
	for i := range nodes {
		result[i] = nodes[i]
	}
	return discovery.Result{
		Cacheable: true,
		CacheKey:  desc,
		Instances: result,
	}, nil
}

func (e *etcdResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

func (e *etcdResolver) Name() string {
	return EtcdResolver
}
