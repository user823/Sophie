package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/vo"
	"github.com/user823/Sophie/pkg/utils"
)

type SysMenu struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	MenuId         int64      `json:"menuId,omitempty" gorm:"column:menu_id" query:"menuId"`
	MenuName       string     `json:"menuName,omitempty" gorm:"column:menu_name" query:"menuName"`
	ParentName     string     `json:"parentName,omitempty" gorm:"column:parent_name" query:"parentName"`
	ParentId       int64      `json:"parentId,omitempty" gorm:"column:parent_id" query:"parentId"`
	OrderNum       int64      `json:"orderNum,omitempty" gorm:"column:order_num" query:"orderNum"`
	Path           string     `json:"path,omitempty" gorm:"column:path" query:"path"`
	Component      string     `json:"component,omitempty" gorm:"column:component" query:"component"`
	Query          string     `json:"query,omitempty" gorm:"column:query" query:"query"`
	IsFrame        string     `json:"isFrame,omitempty" gorm:"column:is_frame" query:"isFrame"`
	IsCache        string     `json:"isCache,omitempty" gorm:"column:is_cache" query:"isCache"`
	MenuType       string     `json:"menuType,omitempty" gorm:"column:menu_type" query:"menuType"`
	Visible        string     `json:"visible,omitempty" gorm:"column:visible" query:"visible"`
	Status         string     `json:"status,omitempty" gorm:"column:status" query:"status"`
	Perms          string     `json:"perms,omitempty" gorm:"column:perms" query:"perms"`
	Icon           string     `json:"icon,omitempty" gorm:"column:icon" query:"icon"`
	Children       []*SysMenu `json:"children,omitempty" gorm:"-" query:"children"`
}

func (s *SysMenu) TableName() string {
	return "sys_menu"
}

func (s *SysMenu) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysMenu) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

func (s *SysMenu) BuildTreeSelect() vo.TreeSelect {
	children := make([]vo.TreeSelect, 0, len(s.Children))
	for i := range s.Children {
		children = append(children, s.Children[i].BuildTreeSelect())
	}
	return vo.TreeSelect{
		Id:       s.MenuId,
		Label:    s.MenuName,
		Children: children,
	}
}

type MenuList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysMenu `json:"items"`
}
