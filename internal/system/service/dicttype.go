package service

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils/cacheutils"
	"github.com/user823/Sophie/pkg/log"
	"sync"
)

type DictTypeSrv interface {
	// 根据条件分页查询字典类型
	SelectDictTypeList(ctx context.Context, dictType *v1.SysDictType, opts *api.GetOptions) *v1.DictTypeList
	// 根据所有字典类型
	SelectDictTypeAll(ctx context.Context, opts *api.GetOptions) *v1.DictTypeList
	// 根据字典类型查询字典数据
	SelectDictDataByType(ctx context.Context, dictType string, opts *api.GetOptions) *v1.DictDataList
	// 根据字典类型ID查询信息
	SelectDictTypeById(ctx context.Context, dictId int64, opts *api.GetOptions) *v1.SysDictType
	// 根据字典类型查询信息
	SelectDictTypeByType(ctx context.Context, dictType string, opts *api.GetOptions) *v1.SysDictType
	// 批量删除字典信息
	DeleteDictTypeByIds(ctx context.Context, dictIds []int64, opts *api.DeleteOptions) error
	// 加载字典缓存数据
	LoadingDictCache(ctx context.Context)
	// 清空字典缓存数据
	ClearDictCache(ctx context.Context)
	// 重置字典缓存数据
	ResetDictCache(ctx context.Context)
	// 新增保存字典类型信息
	InsertDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.CreateOptions) error
	// 修改保存字典类型信息
	UpdateDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.UpdateOptions) error
	// 校验字典类型是否唯一
	CheckDictTypeUnique(ctx context.Context, dictType *v1.SysDictType, opts *api.GetOptions) bool
}

type dictTypeService struct {
	store store.Factory
}

var _ DictTypeSrv = &dictTypeService{}

var loadDictDataOnce sync.Once

func NewDictTypes(s store.Factory) DictTypeSrv {
	loadDictDataOnce.Do(func() {
		cacheutils.LoadingDictCache(s)
	})
	return &dictTypeService{s}
}

func (s *dictTypeService) SelectDictTypeList(ctx context.Context, dictType *v1.SysDictType, opts *api.GetOptions) *v1.DictTypeList {
	result, total, err := s.store.DictTypes().SelectDictTypeList(ctx, dictType, opts)
	if err != nil {
		return &v1.DictTypeList{ListMeta: api.ListMeta{0}}
	}
	return &v1.DictTypeList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *dictTypeService) SelectDictTypeAll(ctx context.Context, opts *api.GetOptions) *v1.DictTypeList {
	result, total, err := s.store.DictTypes().SelectDictTypeAll(ctx, opts)
	if err != nil {
		return &v1.DictTypeList{ListMeta: api.ListMeta{0}}
	}
	return &v1.DictTypeList{
		ListMeta: api.ListMeta{total},
		Items:    result,
	}
}

func (s *dictTypeService) SelectDictDataByType(ctx context.Context, dictType string, opts *api.GetOptions) *v1.DictDataList {
	dictDatas := cacheutils.GetDictCache(dictType)
	if len(dictDatas) > 0 {
		return &v1.DictDataList{ListMeta: api.ListMeta{int64(len(dictDatas))}, Items: dictDatas}
	}
	dictDatas, total, _ := s.store.DictData().SelectDictDataByType(ctx, dictType, opts)
	if len(dictDatas) > 0 {
		cacheutils.SetDictCache(dictType, dictDatas)
	}
	return &v1.DictDataList{
		ListMeta: api.ListMeta{total},
		Items:    dictDatas,
	}
}
func (s *dictTypeService) SelectDictTypeById(ctx context.Context, dictId int64, opts *api.GetOptions) *v1.SysDictType {
	result, err := s.store.DictTypes().SelectDictTypeById(ctx, dictId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *dictTypeService) SelectDictTypeByType(ctx context.Context, dictType string, opts *api.GetOptions) *v1.SysDictType {
	result, err := s.store.DictTypes().SelectDictTypeByType(ctx, dictType, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *dictTypeService) DeleteDictTypeByIds(ctx context.Context, dictIds []int64, opts *api.DeleteOptions) error {
	for i := range dictIds {
		dictType := s.SelectDictTypeById(ctx, dictIds[i], &api.GetOptions{Cache: true})
		if s.store.DictData().CountDictDataByType(ctx, dictType.DictType, &api.GetOptions{Cache: true}) > 0 {
			return fmt.Errorf("%s 已分配，不能删除", dictType.DictName)
		}
		cacheutils.RemoveDictCache(dictType.DictType)
		return s.store.DictTypes().DeleteDictTypeById(ctx, dictIds[i], opts)
	}
	return nil
}

func (s *dictTypeService) LoadingDictCache(ctx context.Context) {
	cacheutils.LoadingDictCache(s.store)
}

func (s *dictTypeService) ClearDictCache(ctx context.Context) {
	cacheutils.CleanDictCache()
}

func (s *dictTypeService) ResetDictCache(ctx context.Context) {
	cacheutils.CleanDictCache()
	cacheutils.LoadingDictCache(s.store)
}

func (s *dictTypeService) InsertDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.CreateOptions) error {
	// 首先验证格式
	if opts.Validate {
		if err := dictType.Validate(); err != nil {
			return err
		}
	}

	return s.store.DictTypes().InsertDictType(ctx, dictType, opts)
}

func (s *dictTypeService) UpdateDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.UpdateOptions) error {
	oldDict, err := s.store.DictTypes().SelectDictTypeById(ctx, dictType.DictId, &api.GetOptions{Cache: true})
	if err != nil {
		return s.InsertDictType(ctx, dictType, &api.CreateOptions{})
	}
	tx := s.store.Begin()
	if err = tx.DictData().UpdateDictDataType(ctx, oldDict.DictType, dictType.DictType, opts); err != nil {
		log.Infof("%v", oldDict.DictType)
		log.Infof("%s", err.Error())
		tx.Rollback()
		return err
	}
	if err = tx.DictTypes().UpdateDictType(ctx, dictType, opts); err != nil {
		log.Infof("%s", err.Error())
		tx.Rollback()
		return err
	}
	cacheutils.RemoveDictCache(dictType.DictType)
	return tx.Commit()
}

func (s *dictTypeService) CheckDictTypeUnique(ctx context.Context, dictType *v1.SysDictType, opts *api.GetOptions) bool {
	info := s.store.DictTypes().CheckDictTypeUnique(ctx, dictType.DictType, opts)
	if info != nil && info.DictId != dictType.DictId {
		return false
	}
	return true
}
