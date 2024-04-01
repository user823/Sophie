package loadbalance

import (
	"context"
	"github.com/user823/Sophie/internal/schedule/models"
)

// 根据负载均衡策略获取下一个运行节点实例
type Picker interface {
	Next(ctx context.Context, request any) models.Instance
}

// 基于策略模式的负载均衡实现
type LoadBalancer interface {
	GetPicker(pool models.NodePool) Picker
	Name() string
}

type EmptyPicker struct{}

var _ Picker = &EmptyPicker{}

func (e *EmptyPicker) Next(ctx context.Context, request any) models.Instance {
	return nil
}
