package domain

import "github.com/user823/Sophie/api"

type SysNotice struct {
	api.ObjectMeta
	NoticeId      int64  `json:"noticeId" gorm:"column:notice_id"`
	NoticeTitle   string `json:"noticeTitle" gorm:"column:notice_title"`
	NoticeType    string `json:"noticeType" gorm:"column:notice_type"`
	NoticeContent string `json:"noticeContent" gorm:"column:notice_content"`
	Status        string `json:"status" gorm:"column:status"`
}

func (s *SysNotice) TableName() string {
	return "sys_notice"
}

type NoticeList struct {
	api.ListMeta `json:",inline"`
	Items        []SysNotice `json:"items"`
}
