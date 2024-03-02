package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysPost struct {
	api.ObjectMeta `json:",inline,omitempty"`
	PostId         int64  `json:"postId,omitempty" gorm:"column:post_id" query:"postId"`
	PostCode       string `json:"postCode,omitempty" gorm:"column:post_code" query:"postCode"`
	PostName       string `json:"postName,omitempty" gorm:"column:post_name" query:"postName"`
	PostSort       int64  `json:"postSort,omitempty" gorm:"column:post_sort" query:"postSort"`
	Status         string `json:"status,omitempty" gorm:"column:status" query:"status"`
	// 用户是否存在此岗位标识（默认不存在）
	Flag bool `json:"flag" gorm:"-" query:"flag"`
}

func (s *SysPost) TableName() string {
	return "sys_post"
}

func (s *SysPost) String() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysPost) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type PostList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysPost `json:"items"`
}
