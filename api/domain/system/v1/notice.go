package v1

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/pkg/utils"
)

type SysNotice struct {
	api.ObjectMeta `json:",inline,omitempty"`
	NoticeId       int64  `json:"noticeId,omitempty" gorm:"column:notice_id" query:"noticeId"`
	NoticeTitle    string `json:"noticeTitle,omitempty" gorm:"column:notice_title" query:"noticeTitle"`
	NoticeType     string `json:"noticeType,omitempty" gorm:"column:notice_type" query:"noticeType"`
	NoticeContent  string `json:"noticeContent,omitempty" gorm:"column:notice_content" query:"noticeContent"`
	Status         string `json:"status,omitempty" gorm:"column:status" query:"status"`
}

func (s *SysNotice) TableName() string {
	return "sys_notice"
}

func (s *SysNotice) Marshal() string {
	data, _ := jsoniter.Marshal(s)
	return utils.B2s(data)
}

func (s *SysNotice) Unmarshal(str string) {
	data := utils.S2b(str)
	jsoniter.Unmarshal(data, s)
}

type NoticeList struct {
	api.ListMeta `json:",inline"`
	Items        []*SysNotice `json:"items"`
}
