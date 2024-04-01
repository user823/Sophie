package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
)

type ConfigSrv interface {
	// 查询参数配置信息
	SelectConfigById(ctx context.Context, configId int64, opts *api.GetOptions) *v1.SysConfig
	// 根据键名查询参数配置信息
	SelectConfigByKey(ctx context.Context, configKey string, opts *api.GetOptions) string
	// 查询参数配置列表
	SelectConfigList(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) *v1.ConfigList
	// 新增参数配置
	InsertConfig(ctx context.Context, config *v1.SysConfig, opts *api.CreateOptions) error
	// 修改参数配置
	UpdateConfig(ctx context.Context, config *v1.SysConfig, opts *api.UpdateOptions) error
	// 批量删除参数信息
	DeleteConfigByIds(ctx context.Context, configIds []int64, opts *api.DeleteOptions) error
	// 重置参数缓存数据
	ResetConfigCache()
	// 校验参数键名是否唯一
	CheckConfigKeyUnique(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) bool
}

type configService struct {
	store store.Factory
}

func NewConfigs(s store.Factory) ConfigSrv {
	return &configService{store: s}
}

func (s *configService) SelectConfigById(ctx context.Context, configId int64, opts *api.GetOptions) *v1.SysConfig {
	result, err := s.store.Configs().SelectConfigById(ctx, configId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *configService) SelectConfigByKey(ctx context.Context, configKey string, opts *api.GetOptions) string {
	result, err := s.store.Configs().SelectConfig(ctx, &v1.SysConfig{ConfigKey: configKey}, opts)
	if err != nil {
		return ""
	}
	return result.ConfigValue
}

func (s *configService) SelectConfigList(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) *v1.ConfigList {
	result, total, err := s.store.Configs().SelectConfigList(ctx, config, opts)
	if err != nil {
		return &v1.ConfigList{
			ListMeta: api.ListMeta{0},
		}
	}
	return &v1.ConfigList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *configService) InsertConfig(ctx context.Context, config *v1.SysConfig, opts *api.CreateOptions) error {
	return s.store.Configs().InsertConfig(ctx, config, opts)
}

func (s *configService) UpdateConfig(ctx context.Context, config *v1.SysConfig, opts *api.UpdateOptions) error {
	return s.store.Configs().UpdateConfig(ctx, config, opts)
}

func (s *configService) DeleteConfigByIds(ctx context.Context, configIds []int64, opts *api.DeleteOptions) error {
	return s.store.Configs().DeleteConfigByIds(ctx, configIds, opts)
}

func (s *configService) ResetConfigCache() {
	// 使用标准缓存策略，不允许直接操作缓存
	// do nothing
}

func (s *configService) CheckConfigKeyUnique(ctx context.Context, config *v1.SysConfig, opts *api.GetOptions) bool {
	info := s.store.Configs().CheckConfigKeyUnique(ctx, config.ConfigKey, opts)
	if info != nil && info.ConfigId != config.ConfigId {
		return false
	}
	return true
}
