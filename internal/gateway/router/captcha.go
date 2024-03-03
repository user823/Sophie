package router

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/captcha"
	"github.com/user823/Sophie/pkg/utils/strutil"
)

type codeValid struct {
	Code string `json:"code" vd:"required"`
	UUID string `json:"uuid" vd:"required"`
}

type CaptchaController struct {
	captchaEnabled bool
}

func NewCaptchaController(on bool) *CaptchaController {
	return &CaptchaController{on}
}

// 创建验证码
func (f *CaptchaController) CreateCaptcha(ctx context.Context, c *app.RequestContext) {
	rdb := kv.NewKVStore("redis").(kv.RedisStore)
	rdb.SetKeyPrefix(kv.CAPTHA_CODE_KEY)

	// 验证码信息
	captchaGenerator := captcha.NewCaptchaGenerator(captcha.DefaultCaptchaType)
	uuid, captchaInfo, ans := captchaGenerator.Generate()
	log.Infof("uuid: %s, ans: %s", uuid, ans)

	// 测试
	body := map[string]interface{}{"captchaEnabled": f.captchaEnabled, "img": captchaInfo, "uuid": uuid}

	if err := rdb.SetKey(ctx, uuid, ans, utils.SecondToNano(kv.CAPTHA_CODE_KEY_VALID)); err != nil {
		body["captchaEnabled"] = false
		c.JSON(code.SUCCESS, body)
	}

	c.JSON(code.SUCCESS, body)
}

// 校验验证码
func (f *CaptchaController) CheckCaptcha(ctx context.Context, c *app.RequestContext) {
	// 不开启验证码
	if !f.captchaEnabled {
		c.Next(ctx)
		return
	}

	var req codeValid
	if err := c.BindAndValidate(&req); err != nil {
		log.Errorf("Code request bind or validate error: %s", err.Error())
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
	captchaInfo, err := rdb.GetAndDelete(ctx, req.UUID)
	if err != nil || !strutil.EqualIgnoreCase(req.Code, captchaInfo) {
		core.Fail(c, "验证码错误", nil)
		c.Abort()
		return
	}
	c.Next(ctx)
}
