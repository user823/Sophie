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

type mysqlPostStore struct {
	db *gorm.DB
}

var _ store.PostStore = &mysqlPostStore{}

func selectPostVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_post")
}

func (s *mysqlPostStore) SelectPostList(ctx context.Context, post *v1.SysPost, opts *api.GetOptions) ([]*v1.SysPost, int64, error) {
	query := selectPostVo(s.db)
	if post.PostCode != "" {
		query = query.Where("post_code like ?", "%"+post.PostCode+"%")
	}
	if post.Status != "" {
		query = query.Where("status = ? ", post.Status)
	}
	if post.PostName != "" {
		query = query.Where("post_name like ?", "%"+post.PostName+"%")
	}
	query = opts.SQLCondition(query, "")

	var result []*v1.SysPost
	err := query.Find(&result).Error
	return result, utils.CountQuery(query, opts, ""), err
}

func (s *mysqlPostStore) SelectPostAll(ctx context.Context, opts *api.GetOptions) ([]*v1.SysPost, int64, error) {
	return s.SelectPostList(ctx, &v1.SysPost{}, opts)
}

func (s *mysqlPostStore) SelectPostById(ctx context.Context, postid int64, opts *api.GetOptions) (*v1.SysPost, error) {
	query := selectPostVo(s.db).Where("post_id = ?", postid)
	query = opts.SQLCondition(query, "")

	var result v1.SysPost
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(postid, "", ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlPostStore) SelectPostListByUserId(ctx context.Context, userid int64, opts *api.GetOptions) ([]int64, error) {
	query := s.db.Table("sys_post p").Select("p.post_id").Joins(""+
		"left join sys_user_post up on up.post_id = p.post_id").Joins(""+
		"left join sys_user u on u.user_id = up.user_id").Where("u.user_id = ?", userid)
	query = opts.SQLCondition(query, "")

	var ids []int64
	err := query.Find(&ids).Error
	return ids, err
}

func (s *mysqlPostStore) SelectPostsByUserName(ctx context.Context, name string, opts *api.GetOptions) ([]*v1.SysPost, error) {
	query := s.db.Table("sys_post p").Joins(""+
		"left join sys_user_post up on up.post_id = p.post_id").Joins(""+
		"left join sys_user u on u.user_id = up.user_id").Where("u.user_name = ?", name)
	query = opts.SQLCondition(query, "")

	var result []*v1.SysPost
	err := query.Find(&result).Error
	return result, err
}

func (s *mysqlPostStore) DeletePostById(ctx context.Context, postid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("post_id = ?", postid).Delete(&v1.SysPost{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(postid, "", ""))
}

func (s *mysqlPostStore) DeletePostByIds(ctx context.Context, postids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("post_id in ?", postids).Delete(&v1.SysPost{})

	cacheKeys := make([]string, 0, len(postids))
	for i := range postids {
		cacheKeys = append(cacheKeys, s.CacheKey(postids[i], "", ""))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return del.Error
}

func (s *mysqlPostStore) UpdatePost(ctx context.Context, post *v1.SysPost, opts *api.UpdateOptions) error {
	if post.PostId == 0 {
		return fmt.Errorf("更新岗位必须指定id")
	}
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Model(post).Where("post_id = ?", post.PostId).Updates(post).Error
	}

	s.CachedDB().CleanCache(ctx)
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(post.PostId, "", ""))
}

func (s *mysqlPostStore) InsertPost(ctx context.Context, post *v1.SysPost, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(post)
	return create.Error
}

func (s *mysqlPostStore) CheckPostNameUnique(ctx context.Context, name string, opts *api.GetOptions) *v1.SysPost {
	query := selectPostVo(s.db).Where("post_name = ?", name)
	query = opts.SQLCondition(query, "")

	var result v1.SysPost
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, "", name), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

func (s *mysqlPostStore) CheckPostCodeUnique(ctx context.Context, postCode string, opts *api.GetOptions) *v1.SysPost {
	query := selectPostVo(s.db).Where("post_code = ?", postCode)
	query = opts.SQLCondition(query, "")

	var result v1.SysPost
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(0, postCode, ""), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil
	}
	return &result
}

var postCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlPostStore) CachedDB() *cache.CachedDB {
	postCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-poststore-")
		rdsCli.SetRandomExp(true)

		postCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, postCache.rdsCache)
}

func (s *mysqlPostStore) CacheKey(postid int64, postcode string, name string) string {
	// postid:postcode:postname
	return fmt.Sprintf("%d:%s:%s", postid, postcode, name)
}
