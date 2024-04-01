package v1

import (
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/vo"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/validators"
)

type SysDept struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	DeptId         int64      `json:"deptId,omitempty" gorm:"column:dept_id" query:"deptId"`
	ParentId       int64      `json:"parentId,omitempty" gorm:"column:parent_id" query:"parentId"`
	Ancestors      string     `json:"ancestors,omitempty" gorm:"column:ancestors" query:"ancestors"`
	DeptName       string     `json:"deptName,omitempty" gorm:"column:dept_name" validate:"required,min=0,max=30" query:"deptName"`
	OrderNum       int64      `json:"orderNum,omitempty" gorm:"column:order_num" validate:"required" query:"orderNum"`
	Leader         string     `json:"leader,omitempty" gorm:"column:leader" query:"leader"`
	Phone          string     `json:"phone,omitempty" gorm:"column:phone" validate:"min=0,max=11" query:"phone"`
	Email          string     `json:"email,omitempty" gorm:"column:email" validate:"email,min=0,max=50" query:"email"`
	Status         string     `json:"status,omitempty" gorm:"column:status" query:"status"`
	DelFlag        string     `json:"delFlag,omitempty" gorm:"column:del_flag" query:"delFlag"`
	ParentName     string     `json:"parentName,omitempty" gorm:"-" query:"parentName"`
	Children       []*SysDept `json:"children,omitempty" gorm:"-" query:"children"`
}

func (d *SysDept) TableName() string {
	return "sys_dept"
}

func (s *SysDept) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysDept) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

func (s *SysDept) Validate() error {
	vd := validators.GetValidatorOr()
	err := vd.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return buildDeptErrMsg(e)
		}
	}
	return nil
}

func buildDeptErrMsg(err validator.FieldError) error {
	switch err.StructNamespace() {
	case "SysDept.DictName":
		return validators.BuildErrMsgHelper(err, "部门名称")
	case "SysDept.OrderNum":
		return validators.BuildErrMsgHelper(err, "显示顺序")
	case "SysDept.Phone":
		return validators.BuildErrMsgHelper(err, "联系号码")
	case "SysDept.Email":
		return validators.BuildErrMsgHelper(err, "部门邮箱")
	}
	return nil
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
