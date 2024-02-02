package domain

import "github.com/user823/Sophie/api"

type SysUserPost struct {
	UserId int64 `json:"userId" gorm:"column:user_id"`
	PostId int64 `json:"postId" gorm:"column:post_id"`
}

func (s *SysUserPost) TableName() string {
	return "sys_user_post"
}

type UserPostList struct {
	api.ListMeta `json:",inline"`
	Items        []SysUserPost `json:"items"`
}
