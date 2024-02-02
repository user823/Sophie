package domain

import "time"

type SysUserOnline struct {
	// 会话唯一标识
	TokenId       string    `json:"tokenid"`
	Username      string    `json:"username"`
	Ipaddr        string    `json:"ipaddr"`
	LoginLocation string    `json:"loginLocation"`
	Browser       string    `json:"browser"`
	Os            string    `json:"os"`
	LoginTime     time.Time `json:"loginTime"`
}
