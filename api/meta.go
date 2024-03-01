package api

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/pkg/utils"
	"gorm.io/gorm"
	"time"
)

// 请求参数可放在扩展字段中
type Extend map[string]interface{}

func (ext Extend) String() string {
	data, _ := jsoniter.Marshal(ext)
	return string(data)
}

func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend

	_ = jsoniter.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}

	return ext
}

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}

type ObjectMeta struct {
	CreateBy     string    `json:"createBy,omitempty" gorm:"column:create_by"`
	CreatedAt    time.Time `json:"createTime,omitempty" gorm:"column:create_time"`
	UpdateBy     string    `json:"updateBy,omitempty" gorm:"column:update_by"`
	UpdatedAt    time.Time `json:"updateTime,omitempty" gorm:"column:update_time"`
	Remark       string    `json:"remark,omitempty" gorm:"column:remark"`
	Extend       Extend    `json:"extend,omitempty" gorm:"-"`
	ExtendShadow string    `json:"-" gorm:"column:extend_shadow"`
}

func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// BeforeUpdate run before update database record.
func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
	if obj.ExtendShadow == "" {
		return nil
	}
	if err := jsoniter.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}

type GetOptions struct {
	// 启用缓存
	Cache bool `json:"cache,omitempty"`
	// 分页信息
	PageNum       int64  `json:"pageNum,omitempty"`
	PageSize      int64  `json:"pageSize,omitempty"`
	OrderByColumn string `json:"orderByColumn,omitempty"`
	IsAsc         bool   `json:"isAsc,omitempty"`
	// 日期范围
	BeginTime int64 `json:"beginTime,omitempty"`
	EndTime   int64 `json:"endTime,omitempty"`
}

func (g *GetOptions) StartPage() {
	if g.PageNum <= 0 {
		g.PageNum = 1
	}

	if g.PageSize <= 0 {
		g.PageSize = 10
	}
}

func (g *GetOptions) SQLCondition(db *gorm.DB, timeRangeColumn string) *gorm.DB {

	if timeRangeColumn != "" {
		if g.BeginTime != 0 {
			db = db.Where(timeRangeColumn+" >= ?", utils.Second2Time(g.BeginTime))
		}
		if g.EndTime != 0 {
			db = db.Where(timeRangeColumn+" <= ?", utils.Second2Time(g.EndTime))
		}
	}

	if g.PageNum > 0 && g.PageSize > 0 {
		offset := (g.PageNum - 1) * g.PageSize
		db = db.Offset(int(offset)).Limit(int(g.PageSize))
	}

	if g.OrderByColumn != "" {
		direction := "ASC"
		if !g.IsAsc {
			direction = "DESC"
		}
		db = db.Order(g.OrderByColumn + " " + direction)
	}

	return db
}

type CreateOptions struct{}

func (c *CreateOptions) SQLCondition(db *gorm.DB) *gorm.DB {
	return db
}

type DeleteOptions struct {
	// 是否物理删除
	Unscoped bool `json:"unscoped"`
}

func (d *DeleteOptions) SQLCondition(db *gorm.DB) *gorm.DB {
	return db
}

type UpdateOptions struct{}

func (d *UpdateOptions) SQLCondition(db *gorm.DB) *gorm.DB {
	return db
}
