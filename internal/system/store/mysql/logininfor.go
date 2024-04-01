package mysql

import (
	"context"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils"
	"gorm.io/gorm"
)

type mysqlLogininforStore struct {
	db *gorm.DB
}

var _ store.LogininforStore = &mysqlLogininforStore{}

func (s *mysqlLogininforStore) InsertLogininfor(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(logininfor)
	return create.Error
}

func (s *mysqlLogininforStore) SelectLogininforList(ctx context.Context, logininfor *v1.SysLogininfor, opts *api.GetOptions) ([]*v1.SysLogininfor, int64, error) {
	query := s.db.Table("sys_logininfor")
	if logininfor.Ipaddr != "" {
		query = query.Where("ipaddr like ?", "%"+logininfor.Ipaddr+"%")
	}
	if logininfor.Status != "" {
		query = query.Where("status = ?", logininfor.Status)
	}
	if logininfor.UserName != "" {
		query = query.Where("user_name like ?", "%"+logininfor.UserName+"%")
	}
	query = opts.SQLCondition(query, "create_time")
	query = query.Order("info_id DESC")

	var result []*v1.SysLogininfor
	err := query.Find(&result).Error
	return result, utils.CountQuery(query, opts, "create_time"), err
}

func (s *mysqlLogininforStore) DeleteLogininforByIds(ctx context.Context, ids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("sys_logininfor.info_id in ?", ids).Delete(&v1.SysLogininfor{})
	return del.Error
}

func (s *mysqlLogininforStore) CleanLogininfor(ctx context.Context, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("true").Delete(&v1.SysLogininfor{})
	return del.Error
}
