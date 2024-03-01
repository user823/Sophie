package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type DictTypeStore interface {
	// 根据条件分页查询字典类型
	SelectDictTypeList(ctx context.Context, dictType *v1.SysDictType, opts *api.GetOptions) ([]*v1.SysDictType, error)
	// 查询所有字典类型
	SelectDictTypeAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysDictType, error)
	// 根据字典类型id查询信息
	SelectDictTypeById(ctx context.Context, dictid int64, opts *api.GetOptions) (*v1.SysDictType, error)
	// 根据字典类型查询信息
	SelectDictTypeByType(ctx context.Context, dictType string, opts *api.GetOptions) (*v1.SysDictType, error)
	// 根据字典id删除字典信息
	DeleteDictTypeById(ctx context.Context, dictid int64, opts *api.DeleteOptions) error
	// 批量删除字典类型信息
	DeleteDictTypeByIds(ctx context.Context, dictids []int64, opts *api.DeleteOptions) error
	// 新增字典类型信息
	InsertDictType(ctx context.Context, dictType *v1.SysDictType, otps *api.CreateOptions) error
	// 修改字典类型信息
	UpdateDictType(ctx context.Context, dictType *v1.SysDictType, opts *api.UpdateOptions) error
	// 校验字典类型名称是否唯一
	CheckDictTypeUnique(ctx context.Context, dictType string, opts *api.GetOptions) *v1.SysDictType
}
