package secure

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gateway/v1"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strings"
)

const (
	REQUIRE_PERMISSION = "requirePermissions"
	ALL_PERMISSIONS    = "*.*.*"
)

func RequirePermissions(permissions string, opts ...Option) app.HandlerFunc {
	options := DefaultOption()
	// 修改选项设置
	for i := range opts {
		opts[i](options)
	}
	return func(ctx context.Context, c *app.RequestContext) {
		// 不需要验证权限
		if permissions == "" {
			next(ctx, c)
			return
		}
		requirePermissions := strings.Split(permissions, options.Separator)

		// 获取登录信息
		data, ok := c.Get(api.LOGIN_INFO_KEY)
		if !ok || data == nil {
			core.WriteResponse(c, core.ErrResponse{Code: code.UNAUTHRIZED, Message: "登录信息失效，请重新登录"})
		}
		loginInfo := data.(*v1.LoginUser)

		// 获取用户权限
		perms := loginInfo.Permissions
		if options.Logic == AND {
			for _, permission := range requirePermissions {
				if !hasPerm(permission, perms) {
					options.Forbidden(ctx, c)
					return
				}
			}
			next(ctx, c)
			return
		} else if options.Logic == OR {
			for _, permission := range requirePermissions {
				if hasPerm(permission, perms) {
					next(ctx, c)
					return
				}
			}
			options.Forbidden(ctx, c)
			return
		} else {
			log.Panicf("unsupported permissions relationship: %d", options.Logic)
			core.WriteResponse(c, core.ErrResponse{Code: code.ERROR, Message: "系统内部错误"})
			return
		}
	}
}

func hasPerm(requirePerm string, perms []string) bool {
	compare := func(str string, target string) bool {
		if target == ALL_PERMISSIONS {
			return true
		}
		return strutil.SimpleMatch(str, target)
	}
	return strutil.CompareAny(compare, requirePerm, perms...)
}

func next(ctx context.Context, c *app.RequestContext) {
	// 权限校验通过则设置登录信息到ctx 中，便于后续调用
	data, _ := c.Get(api.LOGIN_INFO_KEY)
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, data)
	c.Next(ctx)
}
