package loadbalance

import (
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/utils/hash"
)

const (
	ConsistentHash = "consistent_hash"
)

type consistentHashBalancer struct {
	// 加密算法
	hasher hash.Hasher
	// 虚拟节点数目
	replica int
}

func NewConsistentHashBalancer(hasher hash.Hasher, replica int) LoadBalancer {
	return &consistentHashBalancer{hasher, replica}
}

func (c *consistentHashBalancer) Name() string {
	return "consistent_hash_balancer"
}

func (c *consistentHashBalancer) GetPicker(pool models.NodePool) Picker {
	nodes := pool.OnlineNodes()
	instances := make([]models.Instance, len(nodes))
	for i := range nodes {
		instances[i] = nodes[i]
	}
	return newConsistentHashPicker(instances, c.hasher, c.replica)
}
