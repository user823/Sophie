package domain

import "github.com/user823/Sophie/api"

type SysPost struct {
	api.ObjectMeta
	PostId   int64  `json:"postId" gorm:"column:post_id"`
	PostCode string `json:"postCode" gorm:"column:post_code"`
	PostName string `json:"postName" gorm:"column:post_name"`
	PostSort int64  `json:"postSort" gorm:"column:post_sort"`
	Status   string `json:"status" gorm:"column:status"`
	Flag     bool   `json:"flag" gorm:"column:flag"`
}

func (s *SysPost) TableName() string {
	return "sys_post"
}

type PostList struct {
	api.ListMeta `json:",inline"`
	Items        []SysPost `json:"items"`
}
