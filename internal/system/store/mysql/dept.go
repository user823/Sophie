package mysql

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/cache"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/pkg/db/kv"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type mysqlDeptStore struct {
	db *gorm.DB
}

var _ store.DeptStore = &mysqlDeptStore{}

func selectDeptVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_dept d")
}

func (s *mysqlDeptStore) SelectDeptList(ctx context.Context, dept *v1.SysDept, opts *api.GetOptions) ([]*v1.SysDept, error) {
	query := selectDeptVo(s.db).Where("del_flag = 0")
	if dept.DeptId != 0 {
		query = query.Where("dept_id = ?", dept.DeptId)
	}
	if dept.ParentId != 0 {
		query = query.Where("parent_id = ?", dept.ParentId)
	}
	if dept.DeptName != "" {
		query = query.Where("dept_name like ?", "%"+dept.DeptName+"%")
	}
	if dept.Status != "" {
		query = query.Where("status = ?", dept.Status)
	}
	query = opts.SQLCondition(query, "")
	query, err := dateScopeFromCtx(ctx, query, "", "d")
	if err != nil {
		return []*v1.SysDept{}, err
	}

	var result []*v1.SysDept
	err = query.Find(&result).Error
	return result, err
}

func (s *mysqlDeptStore) SelectDeptListByRoleId(ctx context.Context, roleid int64, deptCheckStrictly bool, opts *api.GetOptions) ([]int64, error) {
	query := s.db.Table("sys_dept d").Select("d.dept_id").Joins(""+
		"left join sys_role_dept rd on d.dept_id = rd.dept_id").Where("rd.role_id = ?", roleid)
	if deptCheckStrictly {
		query = query.Where("d.dept_id not in (select d.parent_id from sys_dept d join sys_role_dept rd on d.dept_id = rd.dept_id and rd.role_id = ?)", roleid)
	}
	query = query.Order("d.parent_id ASC, d.order_num ASC")
	query = opts.SQLCondition(query, "")

	var result []int64
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlDeptStore) SelectDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) (*v1.SysDept, error) {
	query := selectDeptVo(s.db).Where("dept_id = ?", deptId)
	query = opts.SQLCondition(query, "")

	var result v1.SysDept
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(deptId, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlDeptStore) SelectChildrenDeptById(ctx context.Context, deptid int64, opts *api.GetOptions) ([]*v1.SysDept, error) {
	query := s.db.Table("sys_dept").Where("find_in_set(?, ancestors)", deptid)
	query = opts.SQLCondition(query, "")

	var result []*v1.SysDept
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlDeptStore) SelectNormalChildrenDeptById(ctx context.Context, deptId int64, opts *api.GetOptions) int {
	query := s.db.Table("sys_dept d").Where("d.status = 0 and d.del_flag = 0 and find_in_set(?, ancestors)", deptId)
	query = opts.SQLCondition(query, "")
	var result int64
	query.Count(&result)
	return int(result)
}

func (s *mysqlDeptStore) HasChildByDeptId(ctx context.Context, deptid int64, opts *api.GetOptions) bool {
	query := s.db.Table("sys_dept d").Where("d.del_flag = 0 and d.parent_id = ?", deptid)
	query = opts.SQLCondition(query, "")
	var result int64
	query.Count(&result)
	return result > 0
}

func (s *mysqlDeptStore) CheckDeptExistUser(ctx context.Context, deptid int64, opts *api.GetOptions) bool {
	query := s.db.Table("sys_user u").Where("u.dept_id = ? and u.del_flag = 0", deptid)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return result > 0
}

func (s *mysqlDeptStore) CheckDeptNameUnique(ctx context.Context, name string, deptid int64, opts *api.GetOptions) *v1.SysDept {
	query := selectDeptVo(s.db).Where("dept_name = ? and parent_id = ? and del_flag = 0", name, deptid)
	query = opts.SQLCondition(query, "")

	var result v1.SysDept
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(&result).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, name), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

func (s *mysqlDeptStore) InsertDept(ctx context.Context, dept *v1.SysDept, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(dept)
	return create.Error
}

func (s *mysqlDeptStore) UpdateDept(ctx context.Context, dept *v1.SysDept, opts *api.UpdateOptions) error {
	if dept.DeptId == 0 {
		return fmt.Errorf("更新部门必须指定id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(dept).Where("dept_id = ?", dept.DeptId).Updates(dept).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(dept.DeptId, ""))
}

func (s *mysqlDeptStore) UpdateDeptStatusNormal(ctx context.Context, deptids []int64, opts *api.UpdateOptions) error {
	update := opts.SQLCondition(s.db).Model(&v1.SysDept{}).Where("dept_id in ?", deptids).Update("status", "0")

	cacheKeys := make([]string, 0, len(deptids))
	for i := range deptids {
		cacheKeys = append(cacheKeys, s.CacheKey(deptids[i], ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return update.Error
}

func (s *mysqlDeptStore) UpdateDeptChildren(ctx context.Context, depts []*v1.SysDept, opts *api.UpdateOptions) error {
	cacheKeys := make([]string, 0, len(depts))
	for i := range depts {
		cacheKeys = append(cacheKeys, s.CacheKey(depts[i].DeptId, ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)

	var builder strings.Builder
	ids := make([]int64, len(depts))
	builder.WriteString("case dept_id")
	for i := range depts {
		builder.WriteString(fmt.Sprintf(" when %d then %s", depts[i].DeptId, depts[i].Ancestors))
		ids = append(ids, depts[i].DeptId)
	}
	builder.WriteString(" end")
	query := builder.String()

	sql := "UPDATE sys_dept SET ancestors = " + query + " WHERE dept_id in ?"
	return s.db.Exec(sql, ids).Error
}

func (s *mysqlDeptStore) DeleteDeptById(ctx context.Context, deptid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("dept_id = ?", deptid).Delete(&v1.SysDept{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(deptid, ""))
}

var deptCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlDeptStore) CachedDB() *cache.CachedDB {
	deptCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis").(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-deptstore-")
		rdsCli.SetRandomExp(true)

		deptCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, deptCache.rdsCache)
}

func (s *mysqlDeptStore) CacheKey(dictDataId int64, name string) string {
	// id:name
	return fmt.Sprintf("%d:%s", dictDataId, name)
}
