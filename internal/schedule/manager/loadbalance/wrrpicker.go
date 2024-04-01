package loadbalance

import (
	"context"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/utils/randutil"
)

// 基于权重的轮训负载均衡
type WeightedRoundRobinPicker struct {
	instances     []models.Instance
	currentWeight []int
}

func newWeightedRoundRobinPicker(instances []models.Instance) Picker {
	picker := new(WeightedRoundRobinPicker)
	picker.instances = instances
	picker.currentWeight = make([]int, len(instances))
	return picker
}

func (w *WeightedRoundRobinPicker) Next(ctx context.Context, req any) models.Instance {
	if len(w.instances) == 0 {
		return nil
	}

	// 记录多个最大的权重
	var max_queue []int
	// 更新前恢复权重
	selectedNode := 0
	total := 0

	for i, instance := range w.instances {
		w.currentWeight[i] += instance.Weight()
		total += w.currentWeight[i]
		if w.currentWeight[i] > w.currentWeight[selectedNode] {
			selectedNode = i
			max_queue = []int{i}
		} else if w.currentWeight[i] == w.currentWeight[selectedNode] {
			max_queue = append(max_queue, i)
		}
	}
	if len(max_queue) > 1 {
		selectedNode = max_queue[randutil.Intn(len(max_queue))]
	}

	// 选取后减去权重
	w.currentWeight[selectedNode] -= total
	return w.instances[selectedNode]
}
