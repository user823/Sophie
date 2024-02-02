package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	guuid "github.com/google/uuid"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

type codeValid struct {
	Code string `json:"code" form:"code" cookie:"code" vd:"required"`
	UUID string `json:"uuid" form:"uuid" cookie:"uuid" vd:"required"`
}

// 创建验证码
func createCaptcha(ctx context.Context, c *app.RequestContext) {
	rdb := kv.NewKVStore("redis").(kv.RedisStore)
	rdb.SetKeyPrefix(kv.CAPTHA_CODE_KEY)

	// 验证码信息
	uuid := guuid.NewString()

}

// 校验验证码
func checkCaptcha(ctx context.Context, c *app.RequestContext) {
	var req codeValid
	if err := c.BindAndValidate(&req); err != nil {
		core.Fail(c, "无效验证码", nil)
		c.Abort()
		return
	}
	rdb := kv.NewKVStore("redis").(kv.RedisStore)
	rdb.SetKeyPrefix(kv.CAPTHA_CODE_KEY)
	if !rdb.Connected() {
		core.WriteResponse(c, core.ErrResponse{Code: code.ERROR, Message: "服务器内部错误，请重试"})
		c.Abort()
		return
	}
	captcha, err := rdb.GetAndDelete(ctx, req.UUID)
	if err != nil || strutil.EqualIgnoreCase(req.Code, captcha) {
		core.Fail(c, "验证码错误", nil)
		c.Abort()
		return
	}
	c.Next(ctx)
}
