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
	"sync"
)

type mysqlMenuStore struct {
	db *gorm.DB
}

var _ store.MenuStore = &mysqlMenuStore{}

func selectMenuVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_menu")
}

func (s *mysqlMenuStore) SelectMenuList(ctx context.Context, menu *v1.SysMenu, opts *api.GetOptions) ([]*v1.SysMenu, error) {
	query := selectMenuVo(s.db)
	if menu.MenuName != "" {
		query = query.Where("menu_name like ?", "%"+menu.MenuName+"%")
	}
	if menu.Visible != "" {
		query = query.Where("visible = ?", menu.Visible)
	}
	if menu.Status != "" {
		query = query.Where("status = ?", menu.Status)
	}
	query = opts.SQLCondition(query, "")
	query = query.Order("parent_id, order_num")

	var result []*v1.SysMenu
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuPerms(ctx context.Context, opts *api.GetOptions) ([]string, error) {
	query := s.db.Table("sys_menu m").Select("DISTINCT m.perms").Joins("" +
		"left join sys_role_menu rm on m.menu_id = rm.menu_id").Joins("" +
		"left join sys_user_role ur on rm.role_id = ur.role_id")
	query = opts.SQLCondition(query, "")

	var result []string
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuListByUserId(ctx context.Context, menu *v1.SysMenu, userid int64, opts *api.GetOptions) ([]*v1.SysMenu, error) {
	query := s.db.Table("sys_menu m").Joins(""+
		"left join sys_role_menu rm on m.menu_id = rm.menu_id").Joins(""+
		"left join sys_user_role ur on rm.role_id = ur.role_id").Joins(""+
		"left join sys_role ro on ur.role_id = ro.role_id").Where("ur.user_id = ?", userid).Distinct("" +
		"m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.query, m.visible, m.status, m.perms, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time")
	if menu.MenuName != "" {
		query = query.Where("m.menu_name like ?", "%"+menu.MenuName+"%")
	}
	if menu.Visible != "" {
		query = query.Where("m.visible = ?", menu.Visible)
	}
	if menu.Status != "" {
		query = query.Where("m.status = ? ", menu.Status)
	}
	query = opts.SQLCondition(query, "")
	query = query.Order("m.parent_id, m.order_num")

	var result []*v1.SysMenu
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuPermsByRoleId(ctx context.Context, roleid int64, opts *api.GetOptions) ([]string, error) {
	query := s.db.Table("sys_menu m").Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").Where(""+
		"m.status = 0 and rm.role_id = ?", roleid).Select("DISTINCT m.perms")
	query = opts.SQLCondition(query, "")

	var result []string
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuPermsByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]string, error) {
	query := s.db.Table("sys_menu m").Joins(""+
		"left join sys_role_menu rm on m.menu_id = rm.menu_id").Joins(""+
		"left join sys_user_role ur on rm.role_id = ur.role_id").Joins(""+
		"left join sys_role r on r.role_id = ur.role_id").Where("m.status = 0 and r.status = 0 and ur.user_id = ?", userid).Select("DISTINCT m.perms")
	query = opts.SQLCondition(query, "")

	var result []string
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuTreeAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysMenu, error) {
	query := s.db.Table("sys_menu m").Where("m.menu_type in ('M', 'C') and m.status = 0").Select("" +
		"DISTINCT m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.query, m.visible, m.status, m.perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time")
	query = opts.SQLCondition(query, "")
	query = query.Order("m.parent_id, m.order_num")

	var result []*v1.SysMenu
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuTreeByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]*v1.SysMenu, error) {
	query := s.db.Table("sys_menu m").Joins(""+
		"left join sys_role_menu rm on m.menu_id = rm.menu_id").Joins(""+
		"left join sys_user_role ur on rm.role_id = ur.role_id").Joins(""+
		"left join sys_role ro on ro.role_id = ur.role_id").Where("ur.user_id = ? and m.menu_type in ('M', 'C') and m.status = 0 and ro.status = 0", userid).Select("" +
		"DISTINCT m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.query, m.visible, m.status, m.perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time")
	query = opts.SQLCondition(query, "")
	query = query.Order("m.parent_id, m.order_num")

	var result []*v1.SysMenu
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuListByRoleId(ctx context.Context, roleid int64, menuCheckStrictly bool, opts *api.GetOptions) ([]int64, error) {
	query := s.db.Table("sys_menu m").Joins(""+
		"left join sys_role_menu rm on m.menu_id = rm.menu_id").Where("rm.role_id = ?", roleid).Select("m.menu_id")
	if menuCheckStrictly {
		query = query.Where("m.menu_id not in (select m.parent_id from sys_menu m join sys_role_menu rm on m.menu_id = rm.menu_id and rm.role_id = ?)", roleid)
	}
	query = opts.SQLCondition(query, "")
	query = query.Order("m.parent_id, m.order_num")

	var result []int64
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlMenuStore) SelectMenuById(ctx context.Context, menuid int64, opts *api.GetOptions) (*v1.SysMenu, error) {
	query := selectMenuVo(s.db).Where("menu_id = ?", menuid)
	query = opts.SQLCondition(query, "")

	var result v1.SysMenu
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(menuid, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlMenuStore) HasChildByMenuId(ctx context.Context, menuid int64, opts *api.GetOptions) bool {
	query := s.db.Table("sys_menu m").Where("m.parent_id = ?", menuid)
	query = opts.SQLCondition(query, "")

	var result int64
	query.Count(&result)
	return result > 0
}

func (s *mysqlMenuStore) InsertMenu(ctx context.Context, menu *v1.SysMenu, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(menu)
	return create.Error
}

func (s *mysqlMenuStore) UpdateMenu(ctx context.Context, menu *v1.SysMenu, opts *api.UpdateOptions) error {
	if menu.MenuId == 0 {
		return fmt.Errorf("更新菜单必须指明id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("menu_id = ?", menu.MenuId).Updates(menu).Error
	}

	s.CachedDB().CleanCache(ctx)
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(menu.MenuId, ""))
}

func (s *mysqlMenuStore) DeleteMenuById(ctx context.Context, menuid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("menu_id = ?", menuid).Delete(&v1.SysMenu{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(menuid, ""))
}

func (s *mysqlMenuStore) CheckMenuNameUnique(ctx context.Context, menuName string, parentid int64, opts *api.GetOptions) *v1.SysMenu {
	query := selectMenuVo(s.db).Where("menu_name = ? and parent_id = ?", menuName, parentid)
	query = opts.SQLCondition(query, "")

	var result v1.SysMenu
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, menuName), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

var menuCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlMenuStore) CachedDB() *cache.CachedDB {
	menuCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-menustore-")
		rdsCli.SetRandomExp(true)

		menuCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, menuCache.rdsCache)
}

func (s *mysqlMenuStore) CacheKey(menuId int64, name string) string {
	// id:name
	return fmt.Sprintf("%d:%s", menuId, name)
}
