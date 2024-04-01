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

type mysqlNoticeStore struct {
	db *gorm.DB
}

var _ store.NoticeStore = &mysqlNoticeStore{}

func selectNoticeVo(db *gorm.DB) *gorm.DB {
	return db.Table("sys_notice")
}

func (s *mysqlNoticeStore) SelectNoticeById(ctx context.Context, noticeid int64, opts *api.GetOptions) (*v1.SysNotice, error) {
	query := selectNoticeVo(s.db).Where("notice_id = ?", noticeid)
	query = opts.SQLCondition(query, "")

	var result v1.SysNotice
	var err error
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		return query.First(v).Error
	}

	if opts.Cache {
		err = s.CachedDB().QueryRow(ctx, s.CacheKey(noticeid), &result, queryFn)
	} else {
		err = queryFn(ctx, query, &result)
	}

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *mysqlNoticeStore) SelectNoticeList(ctx context.Context, notice *v1.SysNotice, opts *api.GetOptions) ([]*v1.SysNotice, int64, error) {
	query := selectNoticeVo(s.db)
	if notice.NoticeTitle != "" {
		query = query.Where("notice_title like ?", "%"+notice.NoticeTitle+"%")
	}
	if notice.NoticeType != "" {
		query = query.Where("notice_type = ?", notice.NoticeType)
	}
	if notice.CreateBy != "" {
		query = query.Where("create_by like ?", "%"+notice.CreateBy+"%")
	}
	query = opts.SQLCondition(query, "")

	var result []*v1.SysNotice
	err := query.Find(&result).Error
	return result, utils.CountQuery(query, opts, ""), err
}

func (s *mysqlNoticeStore) InsertNotice(ctx context.Context, notice *v1.SysNotice, opts *api.CreateOptions) error {
	create := opts.SQLCondition(s.db).Create(notice)
	return create.Error
}

func (s *mysqlNoticeStore) UpdateNotice(ctx context.Context, notice *v1.SysNotice, opts *api.UpdateOptions) error {
	if notice.NoticeId == 0 {
		return fmt.Errorf("更新公告必须指定id")
	}

	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("notice_id = ?", notice.NoticeId).Updates(notice).Error
	}

	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(notice.NoticeId))
}

func (s *mysqlNoticeStore) DeleteNoticeById(ctx context.Context, noticeid int64, opts *api.DeleteOptions) error {
	execFn := func(ctx context.Context, db *gorm.DB) error {
		return opts.SQLCondition(s.db).Where("notice_id = ?", noticeid).Delete(&v1.SysNotice{}).Error
	}
	return s.CachedDB().Exec(ctx, execFn, s.CacheKey(noticeid))
}

func (s *mysqlNoticeStore) DeleteNoticeByIds(ctx context.Context, noticeids []int64, opts *api.DeleteOptions) error {
	del := opts.SQLCondition(s.db).Where("notice_id in ?", noticeids).Delete(&v1.SysNotice{})

	cacheKeys := make([]string, 0, len(noticeids))
	for i := range noticeids {
		cacheKeys = append(cacheKeys, s.CacheKey(noticeids[i]))
	}
	s.CachedDB().DelCache(ctx, cacheKeys...)
	return del.Error
}

var noticeCache = struct {
	rdsCache cache.Cache
	once     sync.Once
}{}

func (s *mysqlNoticeStore) CachedDB() *cache.CachedDB {
	noticeCache.once.Do(func() {
		rdsCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
		rdsCli.SetKeyPrefix("sophie-system-noticestore-")
		rdsCli.SetRandomExp(true)

		noticeCache.rdsCache = cache.NewRedisCache(rdsCli, cache.NewSingleFlight())
	})

	return cache.NewCachedDB(s.db, noticeCache.rdsCache)
}

func (s *mysqlNoticeStore) CacheKey(noticeId int64) string {
	// id
	return fmt.Sprintf("%d", noticeId)
}
