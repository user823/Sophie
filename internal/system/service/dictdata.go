package service

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils/cacheutils"
)

type DictDataSrv interface {
	// 根据条件分页查询字典数据
	SelectDictDataList(ctx context.Context, dictData *v1.SysDictData, opts *api.GetOptions) *v1.DictDataList
	// 根据字典类型和字典键值查询字典数据信息
	SelectDictLabel(ctx context.Context, dictType string, dictValue string, opts *api.GetOptions) string
	// 根据字典数据ID查询信息
	SelectDictDataById(ctx context.Context, dictCode int64, opts *api.GetOptions) *v1.SysDictData
	// 批量删除字典数据信息
	DeleteDictDataByIds(ctx context.Context, dictCodes []int64, opts *api.DeleteOptions) error
	// 新增保存字典数据信息
	InsertDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.CreateOptions) error
	// 修改保存字典数据信息
	UpdateDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.UpdateOptions) error
}

type dictDataService struct {
	store store.Factory
}

var _ DictDataSrv = &dictDataService{}

func NewDictDatas(s store.Factory) DictDataSrv {
	return &dictDataService{s}
}

func (s *dictDataService) SelectDictDataList(ctx context.Context, dictData *v1.SysDictData, opts *api.GetOptions) *v1.DictDataList {
	result, total, err := s.store.DictData().SelectDictDataList(ctx, dictData, opts)
	if err != nil {
		return &v1.DictDataList{ListMeta: api.ListMeta{0}}
	}
	return &v1.DictDataList{ListMeta: api.ListMeta{total}, Items: result}
}

func (s *dictDataService) SelectDictLabel(ctx context.Context, dictType string, dictValue string, opts *api.GetOptions) string {
	result, _ := s.store.DictData().SelectDictLabel(ctx, dictType, dictValue, opts)
	return result
}

func (s *dictDataService) SelectDictDataById(ctx context.Context, dictCode int64, opts *api.GetOptions) *v1.SysDictData {
	result, _ := s.store.DictData().SelectDictDataById(ctx, dictCode, opts)
	return result
}

func (s *dictDataService) DeleteDictDataByIds(ctx context.Context, dictCodes []int64, opts *api.DeleteOptions) error {
	for i := range dictCodes {
		data := s.SelectDictDataById(ctx, dictCodes[i], &api.GetOptions{Cache: true})
		cacheutils.RemoveDictCache(data.DictType)
		return s.store.DictData().DeleteDictDataById(ctx, dictCodes[i], opts)
	}
	return nil
}

func (s *dictDataService) InsertDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.CreateOptions) error {
	// 首先验证格式
	if opts.Validate {
		if err := dictData.Validate(); err != nil {
			return err
		}
	}

	return s.store.DictData().InsertDictData(ctx, dictData, opts)
}

func (s *dictDataService) UpdateDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.UpdateOptions) error {
	err := s.store.DictData().UpdateDictData(ctx, dictData, opts)
	if err != nil {
		cacheutils.RemoveDictCache(dictData.DictType)
	}
	return err
}
