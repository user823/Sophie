package mysql

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/cache"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/db/kv"
	"gorm.io/gorm"
	"sync"
)

type mysqlRoleStore struct {
	db *gorm.DB
}

var _ store.RoleStore = &mysqlRoleStore{}

func selectRoleVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_role r").Joins("" +
		"left join sys_user_role ur on ur.role_id = r.role_id").Joins("" +
		"left join sys_user u on u.user_id = ur.user_id").Joins("" +
		"left join sys_dept d on u.dept_id = d.dept_id").Distinct("r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, " +
		"r.menu_check_strictly, r.dept_check_strictly, r.status, r.del_flag, r.create_time, r.create_by, r.remark")
}

func (s *mysqlRoleStore) SelectRoleList(ctx context.Context, role *v1.SysRole, opts *api.GetOptions) ([]*v1.SysRole, int64, error) {
	query := selectRoleVo(s.db).Where("r.del_flag = 0")
	if role.RoleId != 0 {
		query = query.Where("r.role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		query = query.Where("r.role_name = ?", role.RoleName)
	}
	if role.Status != "" {
		query = query.Where("r.status = ?", role.Status)
	}
	if role.RoleKey != "" {
		query = query.Where("r.role_key like ?", "%"+role.RoleKey+"%")
	}
	total := utils.CountQuery(query, opts, "r.create_time")
	query = opts.SQLCondition(query, "r.create_time")
	query, err := dateScopeFromCtx(ctx, query, "", "d")
	if err != nil {
		return []*v1.SysRole{}, 0, err
	}

	var result []*v1.SysRole
	err = query.Find(&result).Error
	return result, total, err
}

func (s *mysqlRoleStore) SelectRolePermissionByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]*v1.SysRole, error) {
	query := selectRoleVo(s.db).Where("r.del_flag = 0 and ur.user_id = ?", userid)
	query = opts.SQLCondition(query, "")

	var result []*v1.SysRole
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlRoleStore) SelectRoleAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysRole, error) {
	query := selectRoleVo(s.db)
	query = opts.SQLCondition(query, "")

	var result []*v1.SysRole
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlRoleStore) SelectRoleListByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]int64, error) {
	query := s.db.Table("sys_role r").Select("r.role_id").Joins(""+
		"left join sys_user_role ur on ur.role_id = r.role_id").Joins(""+
		"left join sys_user u on u.user_id = ur.user_id").Where("u.user_id = ?", userid)
	query = opts.SQLCondition(query, "")

	var result []int64
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlRoleStore) SelectRoleById(ctx context.Context, roleid int64, opts *api.GetOptions) (*v1.SysRole, error) {
	query := selectRoleVo(s.db).Where("r.role_id = ?", roleid)
	query = opts.SQLCondition(query, "")

	var result v1.SysRole
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(roleid, "", ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlRoleStore) SelectRolesByUserName(ctx context.Context, name string, opts *api.GetOptions) ([]*v1.SysRole, error) {
	query := selectRoleVo(s.db).Where("r.del_flag = 0 and u.user_name = ?", name)
	query = opts.SQLCondition(query, "")

	var result []*v1.SysRole
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlRoleStore) CheckRoleNameUnique(ctx context.Context, name string, opts *api.GetOptions) *v1.SysRole {
	query := selectRoleVo(s.db).Where("r.role_name = ? and r.del_flag = 0", name)
	query = opts.SQLCondition(query, "")

	var result v1.SysRole
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, name, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

func (s *mysqlRoleStore) CheckRoleKeyUnique(ctx context.Context, key string, opts *api.GetOptions) *v1.SysRole {
	query := selectRoleVo(s.db).Where("r.role_key = ? and r.del_flag = 0", key)
	query = opts.SQLCondition(query, "")

	var result v1.SysRole
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}
	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, "", key), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result

}

func (s *mysqlRoleStore) UpdateRole(ctx context.Context, role *v1.SysRole, opts *api.UpdateOptions) error {
	if role.RoleId == 0 {
		return fmt.Errorf("更新角色必须指定id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(role).Where("role_id = ?", role.RoleId).Updates(role).Error
	}

	s.CachedDB().CleanCache(ctx)
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(role.RoleId, "", ""))
}

func (s *mysqlRoleStore) InsertRole(ctx context.Context, role *v1.SysRole, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(role)
	return create.Error
}

func (s *mysqlRoleStore) DeleteRoleById(ctx context.Context, roleid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("role_id = ?", roleid).Delete(&v1.SysRole{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(roleid, "", ""))
}

func (s *mysqlRoleStore) DeleteRoleByIds(ctx context.Context, roleids []int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(db).Where("role_id in ?", roleids).Delete(&v1.SysRole{}).Error
	}

	cacheKeys := make([]string, 0, len(roleids))
	for i := range roleids {
		cacheKeys = append(cacheKeys, s.CacheKey(roleids[i], "", ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return execFn(ctx, s.db)
}

var roleCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlRoleStore) CachedDB() *cache.CachedDB {
	roleCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-rolestore-")
		rdsCli.SetRandomExp(true)

		roleCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, roleCache.rdsCache)
}

// 任何更新操作直接删除整个缓存
func (s *mysqlRoleStore) CacheKey(roleId int64, name string, key string) string {
	// roleid:rolename:roleKey
	return fmt.Sprintf("%d:%s:%s", roleId, name, key)
}
