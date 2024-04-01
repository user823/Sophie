package loadbalance

import (
	"context"
	"github.com/user823/Sophie/internal/schedule/models"
	"github.com/user823/Sophie/pkg/utils/randutil"
)

// 基于权重的随机负载均衡策略
type WeightedRandomPicker struct {
	instances []models.Instance
}

func newWeightedRandomPicker(instances []models.Instance) Picker {
	picker := new(WeightedRandomPicker)
	picker.instances = instances
	return picker
}

func (w *WeightedRandomPicker) Next(ctx context.Context, req any) models.Instance {
	if len(w.instances) == 0 {
		return nil
	}
	sameWeight := true
	total := 0
	selectedNode := randutil.Intn(len(w.instances))

	// 计算总权重
	for _, instance := range w.instances {
		total += instance.Weight()
		if instance.Weight() != w.instances[0].Weight() {
			sameWeight = false
		}
	}

	if !sameWeight {
		// 随机从总权重中生成偏移量
		offset := randutil.Intn(total)
		curr := 0
		for i, instance := range w.instances {
			curr += instance.Weight()
			if curr >= offset {
				selectedNode = i
			}
		}
	}
	return w.instances[selectedNode]
}
