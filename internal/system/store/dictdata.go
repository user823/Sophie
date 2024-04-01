package store

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
)

type DictDataStore interface {
	// 根据条件分页查询字典数据
	SelectDictDataList(ctx context.Context, dictData *v1.SysDictData, opts *api.GetOptions) ([]*v1.SysDictData, int64, error)
	// 根据字典类型查询字典数据
	SelectDictDataByType(ctx context.Context, dictType string, opts *api.GetOptions) ([]*v1.SysDictData, int64, error)
	// 根据字典类型和字典键值查询字典数据信息
	SelectDictLabel(ctx context.Context, dictType, dictValue string, opts *api.GetOptions) (string, error)
	// 根据字典数据id查询信息
	SelectDictDataById(ctx context.Context, dictCode int64, opts *api.GetOptions) (*v1.SysDictData, error)
	// 查询字典数据
	CountDictDataByType(ctx context.Context, dictType string, opts *api.GetOptions) int
	// 通过字典id删除字典数据信息
	DeleteDictDataById(ctx context.Context, dictCode int64, opts *api.DeleteOptions) error
	// 批量删除字典数据信息
	DeleteDictDataByIds(ctx context.Context, dictCodes []int64, opts *api.DeleteOptions) error
	// 新增字典数据信息
	InsertDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.CreateOptions) error
	// 修改字典数据信息
	UpdateDictData(ctx context.Context, dictData *v1.SysDictData, opts *api.UpdateOptions) error
	// 同步修改字典类型
	UpdateDictDataType(ctx context.Context, oldType, newType string, opts *api.UpdateOptions) error
}
