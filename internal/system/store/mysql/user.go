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

type mysqlUserStore struct {
	db *gorm.DB
}

var _ store.UserStore = &mysqlUserStore{}

func selectUserVo(db *gorm.DB) *gorm.DB {
	return db.Model(&v1.SysUser{}).Preload("Dept").Preload("Roles").Table("sys_user u")
}

func (s *mysqlUserStore) SelectUserList(ctx context.Context, sysUser *v1.SysUser, opts *api.GetOptions) ([]*v1.SysUser, int64, error) {
	query := selectUserVo(s.db).Where("u.del_flag = 0")
	if sysUser.UserId != 0 {
		query = query.Where("u.user_id = ?", sysUser.UserId)
	}
	if sysUser.Username != "" {
		query = query.Where("u.user_name = ?", sysUser.Username)
	}
	if sysUser.Status != "" {
		query = query.Where("u.status = ?", sysUser.Status)
	}
	if sysUser.Phonenumber != "" {
		query = query.Where("u.phonenumber like ?", "%"+sysUser.Phonenumber+"%")
	}
	if sysUser.DeptId != 0 {
		query = query.Where("u.dept_id = ? OR u.dept_id in (select t.dept_id from sys_dept t where find_in_set(?, ancestors) ))", sysUser.DeptId, sysUser.DeptId)
	}
	query, err := dateScopeFromCtx(ctx, query, "u", "sys_dept")
	total := utils.CountQuery(query, opts, "u.create_time")
	query = opts.SQLCondition(query, "u.create_time")

	if err != nil {
		return []*v1.SysUser{}, 0, err
	}

	var result []*v1.SysUser
	err = query.Find(&result).Error
	return result, total, err
}

// 查询分配了角色的用户列表
func (s *mysqlUserStore) SelectAllocatedList(ctx context.Context, sysUser *v1.SysUser, opts *api.GetOptions) ([]*v1.SysUser, error) {
	query := s.db.Table("sys_user u").Select(""+
		"DISTINCT u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, u.phonenumber, u.status, u.create_time").Joins(""+
		"left join sys_dept d on u.dept_id = d.dept_id").Joins(""+
		"left join sys_user_role ur on u.user_id = ur.user_id").Joins(""+
		"left join sys_role r on r.role_id = ur.role_id").Where("u.del_flag = 0 and r.role_id = ?", sysUser.RoleId)
	if sysUser.Username != "" {
		query = query.Where("u.user_name like ?", "%"+sysUser.Username+"%")
	}
	if sysUser.Phonenumber != "" {
		query = query.Where("u.phonenumber like ?", "%"+sysUser.Phonenumber+"%")
	}
	query = opts.SQLCondition(query, "")
	query, err := dateScopeFromCtx(ctx, query, "u", "sys_dept")
	if err != nil {
		return []*v1.SysUser{}, err
	}

	var result []*v1.SysUser
	err = query.Find(&result).Error
	return result, err
}

// 查询未分配角色的用户列表
func (s *mysqlUserStore) SelectUnallocatedList(ctx context.Context, sysUser *v1.SysUser, opts *api.GetOptions) ([]*v1.SysUser, error) {
	query := s.db.Table("sys_user u").
		Select("DISTINCT u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, u.phonenumber, u.status, u.create_time").
		Joins("LEFT JOIN sys_dept d ON u.dept_id = d.dept_id").
		Joins("LEFT JOIN sys_user_role ur ON u.user_id = ur.user_id").
		Joins("LEFT JOIN sys_role r ON r.role_id = ur.role_id AND (r.role_id != ? OR r.role_id IS NULL)", sysUser.RoleId).
		Where("u.del_flag = 0").
		Not("EXISTS (SELECT 1 FROM sys_user_role WHERE user_id = u.user_id AND role_id = ?)", sysUser.RoleId)
	if sysUser.Username != "" {
		query = query.Where("u.user_name like ?", "%"+sysUser.Username+"%")
	}
	if sysUser.Phonenumber != "" {
		query = query.Where("u.phonenumber like ?", "%"+sysUser.Phonenumber+"%")
	}

	query = opts.SQLCondition(query, "")
	query, err := dateScopeFromCtx(ctx, query, "u", "sys_dept")
	if err != nil {
		return []*v1.SysUser{}, err
	}

	var result []*v1.SysUser
	err = query.Find(&result).Error
	return result, err
}

func (s *mysqlUserStore) SelectUserByUserName(ctx context.Context, name string, opts *api.GetOptions) (*v1.SysUser, error) {
	query := selectUserVo(s.db).Where("u.user_name = ? and u.del_flag = 0", name)
	query = opts.SQLCondition(query, "")

	var result v1.SysUser
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		// cache未命中时默认通过s.db 获取值，因此上面的queryFn 内部使用query查询值
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, name, "", ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *mysqlUserStore) SelectUserById(ctx context.Context, userid int64, opts *api.GetOptions) (*v1.SysUser, error) {
	query := selectUserVo(s.db).Where("u.user_id = ?", userid)
	query = opts.SQLCondition(query, "")

	var result v1.SysUser
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(userid, "", "", ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlUserStore) InsertUser(ctx context.Context, sysUser *v1.SysUser, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(sysUser)
	return create.Error
}

func (s *mysqlUserStore) UpdateUser(ctx context.Context, sysUser *v1.SysUser, opts *api.UpdateOptions) error {
	if sysUser.UserId == 0 {
		return fmt.Errorf("更新用户必须指明id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(sysUser).Where("user_id = ?", sysUser.UserId).Updates(sysUser).Error
	}

	s.CachedDB().CleanCache(ctx)
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(sysUser.UserId, "", "", ""))
}

func (s *mysqlUserStore) UpdateUserAvatar(ctx context.Context, userName, avatar string, opts *api.UpdateOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(&v1.SysUser{}).Where("user_name = ?", userName).Update("avatar", avatar).Error
	}

	s.CachedDB().CleanCache(ctx)
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(0, userName, "", ""))
}

func (s *mysqlUserStore) UpdateUserPwd(ctx context.Context, userName, password string, opts *api.UpdateOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(&v1.SysUser{}).Where("user_name = ?", userName).Update("password", password).Error
	}

	s.CachedDB().CleanCache(ctx)
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(0, userName, "", ""))
}

func (s *mysqlUserStore) DeleteUserById(ctx context.Context, userid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("user_id = ?", userid).Delete(&v1.SysUser{}).Error
	}

	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(userid, "", "", ""))
}

func (s *mysqlUserStore) DeleteUserByIds(ctx context.Context, userids []int64, opts *api.DeleteOptions) error {
	cacheKeys := make([]string, 0, len(userids))
	for i := range userids {
		cacheKeys = append(cacheKeys, s.CacheKey(userids[i], "", "", ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)

	del := opts.SQLCondition(s.db).Where("user_id in ?", userids).Delete(&v1.SysUser{})
	return del.Error
}

// 通过用户名查询id 和 name
func (s *mysqlUserStore) CheckUserNameUnique(ctx context.Context, name string, opts *api.GetOptions) *v1.SysUser {
	query := s.db.Table("sys_user u").Select("u.user_id, u.user_name").Where("u.user_name = ? and u.del_flag = 0 ", name)
	query = opts.SQLCondition(query, "")

	var result v1.SysUser
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, name, "", ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

func (s *mysqlUserStore) CheckPhoneUnique(ctx context.Context, phonenumber string, opts *api.GetOptions) *v1.SysUser {
	query := s.db.Table("sys_user u").Select("u.user_id, u.phonenumber").Where("u.phonenumber = ? and u.del_flag = 0 ", phonenumber)
	query = opts.SQLCondition(query, "")

	var result v1.SysUser
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, "", phonenumber, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, queryFn)
	}

	if err != nil {
		return nil
	}
	return &result
}

func (s *mysqlUserStore) CheckEmailUnique(ctx context.Context, email string, opts *api.GetOptions) *v1.SysUser {
	query := s.db.Table("sys_user u").Select("u.user_id, u.email").Where("u.email = ? and u.del_flag = 0 ", email)
	query = opts.SQLCondition(query, "")

	var result v1.SysUser
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, "", "", email), &result, queryFn)
	} else {
		err = queryFn(ctx, query, queryFn)
	}

	if err != nil {
		return nil
	}
	return &result
}

var userCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlUserStore) CachedDB() *cache.CachedDB {
	userCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-userstore-")
		rdsCli.SetRandomExp(true)

		userCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, userCache.rdsCache)
}

// 任何更新操作都直接删除整个缓存
func (s *mysqlUserStore) CacheKey(userid int64, name string, phone string, email string) string {
	// id:name:phone:email
	return fmt.Sprintf("%d:%s:%s:%s", userid, name, phone, email)
}
