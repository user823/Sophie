package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type ConfigStore interface {
	// 查询参数配置信息
	SelectConfig(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) (*v1.SysConfig, error)
	// 通过id查询配置
	SelectConfigById(ctx context.Context, configid int64, opts *api.GetOptions) (*v1.SysConfig, error)
	// 查询参数配置列表
	SelectConfigList(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) ([]*v1.SysConfig, int64, error)
	// 根据键名查询参数配置信息
	CheckConfigKeyUnique(ctx context.Context, configKey string, opts *api.GetOptions) *v1.SysConfig
	// 新增参数配置
	InsertConfig(ctx context.Context, config *v1.SysConfig, opts *api.CreateOptions) error
	// 修改参数配置
	UpdateConfig(ctx context.Context, config *v1.SysConfig, opts *api.UpdateOptions) error
	// 删除参数配置
	DeleteConfigById(ctx context.Context, configid int64, opts *api.DeleteOptions) error
	// 批量删除参数信息
	DeleteConfigByIds(ctx context.Context, configids []int64, opts *api.DeleteOptions) error
}
