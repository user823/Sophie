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
	SUPER_ADMIN = "admin"
)

func RequireRoles(roles string, opts ...Option) app.HandlerFunc {
	options := DefaultOption()
	// 修改选项设置
	for i := range opts {
		opts[i](options)
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// 没有要求角色
		if roles == "" {
			c.Next(ctx)
			return
		}
		requireRoles := strings.Split(roles, options.Separator)

		// 获取登录信息
		data, ok := c.Get(api.LOGIN_INFO_KEY)
		if !ok || data == nil {
			core.WriteResponse(c, core.ErrResponse{Code: code.UNAUTHRIZED, Message: "登录信息失效，请重新登录"})
		}
		loginInfo := data.(v1.LoginUser)

		// 修改选项设置
		for i := range opts {
			opts[i](options)
		}

		actualRoles := loginInfo.Roles

		if options.Logic == AND {
			for _, role := range requireRoles {
				if !hasRole(role, actualRoles) {
					options.Forbidden(ctx, c)
					return
				}
			}
			c.Next(ctx)
			return
		} else if options.Logic == OR {
			for _, role := range requireRoles {
				if hasRole(role, actualRoles) {
					c.Next(ctx)
					return
				}
			}
			options.Forbidden(ctx, c)
			return
		} else {
			log.Panicf("unsupported role relationship: %d", options.Logic)
			core.WriteResponse(c, core.ErrResponse{Code: code.ERROR, Message: "系统内部错误"})
			return
		}
	}
}

func hasRole(requireRole string, roles []string) bool {
	compare := func(str string, target string) bool {
		if target == SUPER_ADMIN {
			return true
		}
		return strutil.SimpleMatch(str, target)
	}
	return strutil.CompareAny(compare, requireRole, roles...)
}
