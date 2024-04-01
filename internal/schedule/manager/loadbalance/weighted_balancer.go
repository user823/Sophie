package loadbalance

import "github.com/user823/Sophie/internal/schedule/models"

// 基于权重的负载均衡
const (
	WeightedRoundRobin = "weight_round_robin"
	WeightedRandom     = "weight_random"
)

type weightedBalancer struct {
	kind string
}

func NewWeightedBalancer() LoadBalancer {
	return NewWeightedRoundRoinBalancer()
}

func NewWeightedRoundRoinBalancer() LoadBalancer {
	return &weightedBalancer{WeightedRoundRobin}
}

func NewWeightedRandomBalancer() LoadBalancer {
	return &weightedBalancer{WeightedRandom}
}

func (w *weightedBalancer) GetPicker(pool models.NodePool) Picker {
	return w.createPicker(pool)
}

func (w *weightedBalancer) createPicker(pool models.NodePool) (picker Picker) {
	nodes := pool.OnlineNodes()
	instances := make([]models.Instance, len(nodes))
	cnt := 0
	for _, instance := range nodes {
		weight := instance.Weight()
		if weight <= 0 {
			continue
		}
		instances[cnt] = instance
		cnt++
	}
	instances = instances[:cnt]

	switch w.kind {
	case WeightedRoundRobin:
		picker = newWeightedRoundRobinPicker(instances)
	case WeightedRandom:
		picker = newWeightedRandomPicker(instances)
	}
	return
}

func (w *weightedBalancer) Name() string {
	return w.kind
}
