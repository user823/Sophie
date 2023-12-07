package v1

import (
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"time"
)

// Extend used for database extend field
type Extend map[string]interface{}

func (ext Extend) String() string {
	data, _ := jsoniter.Marshal(ext)
	return string(data)
}

func (ext Extend) Merge(extendShadow string) {
	var extend Extend
	_ = jsoniter.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}
}

// ObjectMeta used for Entity common info
type ObjectMeta struct {
	ID           uint64    `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Extend       Extend    `json:"extend,omitempty" gorm:"-" validate:"omitempty"`
	ExtendShadow string    `json:"-" gorm:"column:extend_shadow" validate:"omitempty"`
	Status       string    `json:"status,omitempty" gorm:"status"`
	CreatedBy    string    `json:"created_by,omitempty" gorm:"created_by"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedBy    string    `json:"updated_by,omitempty" gorm:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
}

func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) (err error) {
	obj.ExtendShadow = obj.Extend.String()
	return nil
}

func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) (err error) {
	obj.ExtendShadow = obj.Extend.String()
	return nil
}

func (obj *ObjectMeta) AfterFind(tx *gorm.DB) (err error) {
	err = jsoniter.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend)
	return
}
