package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1/vo"
	"github.com/user823/Sophie/pkg/utils"
)

type SysMenu struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	MenuId         int64      `json:"menuId,omitempty" gorm:"column:menu_id"`
	MenuName       string     `json:"menuName,omitempty" gorm:"column:menu_name"`
	ParentName     string     `json:"parentName,omitempty" gorm:"column:parent_name"`
	ParentId       int64      `json:"parentId,omitempty" gorm:"column:parent_id"`
	OrderNum       int64      `json:"orderNum,omitempty" gorm:"column:order_num"`
	Path           string     `json:"path,omitempty" gorm:"column:path"`
	Component      string     `json:"component,omitempty" gorm:"column:component"`
	Query          string     `json:"query,omitempty" gorm:"column:query"`
	IsFrame        string     `json:"isFrame,omitempty" gorm:"column:is_frame"`
	IsCache        string     `json:"isCache,omitempty" gorm:"column:is_cache"`
	MenuType       string     `json:"menuType,omitempty" gorm:"column:menu_type"`
	Visible        string     `json:"visible,omitempty" gorm:"column:visible"`
	Status         string     `json:"status,omitempty" gorm:"column:status"`
	Perms          string     `json:"perms,omitempty" gorm:"column:perms"`
	Icon           string     `json:"icon,omitempty" gorm:"column:icon"`
	Children       []*SysMenu `json:"-" gorm:"-"`
}

func (s *SysMenu) TableName() string {
	return "sys_menu"
}

func (s *SysMenu) String() string {
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
