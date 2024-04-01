package models

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	v1 "github.com/user823/Sophie/api/domain/schedule/v1"
	"github.com/user823/Sophie/internal/schedule/locker"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNodeNotFound = errors.New("target node not found")
)

// 实现models.NodePool 和 models.JobMap 接口
type EtcdNodePool struct {
	client kv.EtcdStore
}

func NewEtcdNodePool(client kv.EtcdStore) *EtcdNodePool {
	return &EtcdNodePool{client}
}

func (e *EtcdNodePool) sync() (result []Node) {
	ctx, cancel := context.WithTimeout(context.Background(), SyncTimeout*time.Second)
	defer cancel()
	res := e.client.GetKeysAndValuesWithFilter(ctx, NodesPrefix)
	result = make([]Node, 0, len(res))
	cnt := 0
	for _, v := range res {
		var node Node
		if err := jsoniter.Unmarshal(utils.S2b(v), &node); err != nil {
			continue
		}
		result = append(result, node)
		cnt++
	}
	result = result[:cnt]
	return result
}

func (e *EtcdNodePool) Size() int {
	res := e.sync()
	return len(res)
}

func (e *EtcdNodePool) AllNodes() []Node {
	return e.sync()
}

func (e *EtcdNodePool) OnlineNodes() []Node {
	res := e.sync()
	cnt := 0
	for i := range res {
		if res[i].Status == v1.NODE_STATUS_ON {
			res[cnt] = res[i]
			cnt++
		}
	}
	res = res[:cnt]
	return res
}

func (e *EtcdNodePool) GetNode(id string) (Node, error) {
	data, err := e.client.GetKey(context.Background(), ServiceKey(id))
	if err != nil {
		return Node{}, ErrNodeNotFound
	}

	var result Node
	err = jsoniter.Unmarshal(utils.S2b(data), &result)
	if err != nil {
		return Node{}, err
	}

	return result, nil
}

func JobKey(jobid int64) string {
	return fmt.Sprintf("%s:%d", locker.JOB_PREFIX, jobid)
}

func (e *EtcdNodePool) Jobs() map[int64]string {
	if err := locker.RLock(); err != nil {
		log.Error("read lock failed: %s", err.Error())
		return map[int64]string{}
	}
	defer locker.RUnLock()
	ctx, cancel := context.WithTimeout(context.Background(), SyncTimeout*time.Second)
	defer cancel()
	res := e.client.GetKeysAndValuesWithFilter(ctx, locker.JOB_PREFIX)
	result := make(map[int64]string)
	for k, v := range res {
		str := strings.Split(k, ":")
		if len(str) != 2 {
			continue
		}
		jid, err := strconv.ParseInt(str[1], 10, 64)
		if err != nil {
			log.Debugf("get job id failed: %s", err.Error())
			continue
		}
		result[jid] = v
	}
	return result
}

func (e *EtcdNodePool) Exists(jobId int64) (string, bool) {
	jobMp := e.Jobs()
	if a, ok := jobMp[jobId]; ok {
		return a, ok
	}
	return "", false
}

func (e *EtcdNodePool) Create(jobId int64, node string) error {
	if err := locker.Lock(); err != nil {
		return errors.Errorf("lock failed: %s", err.Error())
	}
	defer locker.UnLock()

	return e.client.SetKey(context.Background(), JobKey(jobId), node, 0)
}

func (e *EtcdNodePool) Delete(jobId int64) error {
	if err := locker.Lock(); err != nil {
		return errors.Errorf("lock failed: %s", err.Error())
	}
	defer locker.UnLock()

	if !e.client.DeleteKey(context.Background(), JobKey(jobId)) {
		return errors.Errorf("delete key %s failed", JobKey(jobId))
	}
	return nil
}

func (e *EtcdNodePool) JobsOnNode(nodeid string) (res []int64) {
	mp := e.Jobs()
	for k, v := range mp {
		if v == nodeid {
			res = append(res, k)
		}
	}
	return res
}
