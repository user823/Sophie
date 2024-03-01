package test

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/internal/pkg/cache"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

var (
	rds cache.Cache
)

func initCache() {
	rds = cache.NewRedisCache(client, cache.NewSingleFlight())
}

func TestCache(t *testing.T) {
	v := &TestTable{}

	// set
	err := rds.SetWithExp(context.Background(), "1", &TestTable{Username: "1", Password: "2'"}, 1)
	if err != nil {
		t.Errorf(err.Error())
	}

	// get
	err = rds.Get(context.Background(), "1", v)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(v)

	// take
	err = rds.Take(context.Background(), "3", 10*time.Second, v, func(v any) error {
		v = &TestTable{Username: "2", Password: "3"}
		return nil
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(v)
}

func TestQueryRow(t *testing.T) {
	cachedDB := cache.NewCachedDB(db, rds)

	// 创建数据
	//cachedDB.DB().Create(&TestTable{Username: "11", Password: "11"})
	//cachedDB.DB().Create(&TestTable{Username: "22", Password: "22"})
	//cachedDB.DB().Create(&TestTable{Username: "33", Password: "33"})

	// 查询
	queryFn := func(ctx context.Context, db *gorm.DB, v any) error {
		fmt.Println("执行query")
		return db.First(v, 1).Error
	}

	var result TestTable
	cachedDB.DelCache(context.Background(), "11")
	err := cachedDB.QueryRow(context.Background(), "11", &result, queryFn)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(result)

	// 再次执行QueryRow 希望从缓存中获取值
	var result1 TestTable
	err = cachedDB.QueryRow(context.Background(), "11", &result1, queryFn)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(result1)
}

func TestQueryRowIndex(t *testing.T) {
	var result TestTable
	// cachekey
	key := "c"
	cachedDB := cache.NewCachedDB(db, rds)
	// 主键/唯一索引查询
	indexFn := func(ctx context.Context, db *gorm.DB, v any) (any, error) {
		err := db.WithContext(ctx).Where("username = ?", "33").First(v).Error
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return v.(*TestTable).Id, nil
	}
	// 主键索引查询
	queryFn := func(ctx context.Context, db *gorm.DB, primary any, v any) error {
		return db.Where("id = ?", primary).First(v).Error
	}
	keyer := func(key any) string {
		if key != nil {
			return strconv.Itoa(key.(int))
		}
		return ""
	}
	err := cachedDB.QueryRowIndex(context.Background(), key, &result, keyer, indexFn, queryFn)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(result)
}

func TestCacheSub(t *testing.T) {
	initCache()

	t.Run("test-cache", TestCache)
	t.Run("test-QueryRow", TestQueryRow)
	t.Run("test-QueryRowIndex", TestQueryRowIndex)
}
