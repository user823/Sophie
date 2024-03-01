package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1/vo"
	"github.com/user823/Sophie/pkg/utils"
)

type SysDept struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	DeptId         int64      `json:"deptId,omitempty" gorm:"column:dept_id"`
	ParentId       int64      `json:"parentId,omitempty" gorm:"column:parent_id"`
	Ancestors      string     `json:"ancestors,omitempty" gorm:"column:ancestors"`
	DeptName       string     `json:"deptName,omitempty" gorm:"column:dept_name"`
	OrderNum       int64      `json:"orderNum,omitempty" gorm:"column:order_num"`
	Leader         string     `json:"leader,omitempty" gorm:"column:leader"`
	Phone          string     `json:"phone,omitempty" gorm:"column:phone"`
	Email          string     `json:"email,omitempty" gorm:"column:email"`
	Status         string     `json:"status,omitempty" gorm:"column:status"`
	DelFlag        string     `json:"delFlag,omitempty" gorm:"column:del_flag"`
	ParentName     string     `json:"parentName,omitempty" gorm:"column:parent_name"`
	Children       []*SysDept `json:"-" gorm:"-"`
}

func (d *SysDept) TableName() string {
	return "sys_dept"
}

func (s *SysDept) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysDept) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

func (s *SysDept) BuildTreeSelect() vo.TreeSelect {
	children := make([]vo.TreeSelect, 0, len(s.Children))
	for i := range s.Children {
		children = append(children, s.Children[i].BuildTreeSelect())
	}
	return vo.TreeSelect{
		Id:       s.DeptId,
		Label:    s.DeptName,
		Children: children,
	}
}

type DeptList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysDept `json:"items"`
}
