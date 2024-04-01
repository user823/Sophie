package service

import (
	"context"
	"github.com/user823/Sophie/api"
	sv1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"time"
)

type UserOnlineSrv interface {
	// 通过登录地址查询信息
	// 返回查询结果和总数
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
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)

	expireTime := api.LOGIN_TIMEOUT * time.Second
	_, res := redisCli.GetRollingWindow(ctx, kv.SYS_LOGIN_USER_IDS, expireTime.Milliseconds(), true)
	var keys []string
	for i := range res {
		keys = append(keys, res[i])
	}

	list, err := redisCli.MGetFromHash(ctx, keys)
	if err != nil {
		log.Error("Get login user info error: %s", err.Error())
		return []*sv1.UserOnlineInfo{}
	}

	result := make([]*sv1.UserOnlineInfo, 0, len(list))
	for i := range list {
		var loginUser sv1.Logininfo
		loginUser.Unmarshal(list[i])

		if ipaddr != "" && username != "" {
			if ipaddr == loginUser.Ipaddr && username == loginUser.UserName {
				result = append(result, sv1.LoginInfo2UserOnlineInfo(&loginUser))
			}
		} else if ipaddr != "" {
			if ipaddr == loginUser.Ipaddr {
				result = append(result, sv1.LoginInfo2UserOnlineInfo(&loginUser))
			}
		} else if username != "" {
			if username == loginUser.UserName {
				result = append(result, sv1.LoginInfo2UserOnlineInfo(&loginUser))
			}
		} else {
			result = append(result, sv1.LoginInfo2UserOnlineInfo(&loginUser))
		}

	}
	return result
}

func (s *userOnlineService) ForceLogout(ctx context.Context, tokenId string) {
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.SYS_LOGIN_USER)

	redisCli.DeleteKey(ctx, tokenId)
	redisCli.RemoveFromList(ctx, kv.SYS_LOGIN_USER_IDS, tokenId)
}
