package service

import (
	"context"
	"github.com/user823/Sophie/api"
	sv1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/pkg/db/kv"
	"strconv"
	"strings"
)

type UserOnlineSrv interface {
	// 通过登录地址查询信息
	SelectUserOnline(ctx context.Context, ipaddr string, username string, opts *api.GetOptions) []*sv1.UserOnlineInfo
	// 强制推出
	ForceLogout(ctx context.Context, tokenId string)
}

type userOnlineService struct {
	store store.Factory
}

func NewUserOnlines(s store.Factory) UserOnlineSrv {
	return &userOnlineService{s}
}

func (s *userOnlineService) SelectUserOnline(ctx context.Context, ipaddr string, username string, opts *api.GetOptions) []*sv1.UserOnlineInfo {
	redisCli := kv.NewKVStore("redis").(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)

	keys := redisCli.GetKeys(ctx, "")
	list := make([]*sv1.UserOnlineInfo, 0, len(keys))
	for _, key := range keys {
		logininfo, e := redisCli.GetKey(ctx, key)
		if e != nil {
			continue
		}
		name_token_ip_accesstime := strings.Split(logininfo, ":")
		if ipaddr != "" && username != "" {
			if ipaddr == name_token_ip_accesstime[2] && username == name_token_ip_accesstime[0] {
				list = append(list, generateUserOnlineInfo(name_token_ip_accesstime))
			}
		} else if ipaddr != "" {
			if ipaddr == name_token_ip_accesstime[2] {
				list = append(list, generateUserOnlineInfo(name_token_ip_accesstime))
			}
		} else if username != "" {
			if username == name_token_ip_accesstime[0] {
				list = append(list, generateUserOnlineInfo(name_token_ip_accesstime))
			}
		} else {
			list = append(list, generateUserOnlineInfo(name_token_ip_accesstime))
		}
	}
	return list
}

func generateUserOnlineInfo(name_token_ip_accesstime []string) *sv1.UserOnlineInfo {
	accesstime, _ := strconv.ParseInt(name_token_ip_accesstime[3], 10, 64)
	return &sv1.UserOnlineInfo{
		UserName:  name_token_ip_accesstime[0],
		TokenId:   name_token_ip_accesstime[1],
		Ipaddr:    name_token_ip_accesstime[2],
		LoginTime: accesstime,
	}
}

func (s *userOnlineService) ForceLogout(ctx context.Context, tokenId string) {
	redisCli := kv.NewKVStore("redis").(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)

	redisCli.DeleteKey(ctx, tokenId)
}
