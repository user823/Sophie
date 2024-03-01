package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	jsoniter "github.com/json-iterator/go"
	"github.com/user823/Sophie/api"
	v13 "github.com/user823/Sophie/api/domain/gateway/v1"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/pkg/code"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"net/http"
	"time"
)

// log额外信息key
const (
	TITLE              = "Title"
	BUSINESSTYE        = "BusinessType"
	OperatorType       = "OperatorType"
	IsSaveRequestData  = "IsSaveRequestData"
	IsSaveResponseData = "IsSaveResponseData"
	ExcludParamNames   = "ExcludeParamNames"
)

// log中间件将操作记录到数据库
type LogSaver interface {
	SaveLog(ctx context.Context, operLog *v1.OperLog, options *api.CreateOptions) error
}

type OperLogExtra struct {
	BusinessType       int64
	Title              string
	OperatorType       int64
	IsSaveRequestData  bool
	IsSaveResponseData bool
	ExcludeParamNames  []string
}

// keysAndValues 为环境参数
func Log(logSaver LogSaver, keysAndValues map[string]any) app.HandlerFunc {
	extraInfo := &OperLogExtra{
		BusinessType:       v12.BUSINESSTYPE_OTHER,
		OperatorType:       v12.OPERATORTYPE_MANAGE,
		IsSaveRequestData:  true,
		IsSaveResponseData: true,
		Title:              "",
		ExcludeParamNames:  []string{},
	}
	for k, v := range keysAndValues {
		switch k {
		case TITLE:
			extraInfo.Title = v.(string)
		case BUSINESSTYE:
			extraInfo.BusinessType = v.(int64)
		case OperatorType:
			extraInfo.OperatorType = v.(int64)
		case IsSaveRequestData:
			extraInfo.IsSaveRequestData = v.(bool)
		case IsSaveResponseData:
			extraInfo.IsSaveResponseData = v.(bool)
		case ExcludParamNames:
			extraInfo.ExcludeParamNames = v.([]string)
		}
	}

	return func(ctx context.Context, c *app.RequestContext) {
		operLog := &v1.OperLog{
			CreateTime: utils.Time2Second(time.Now()),
		}
		operLog.Status = v12.BUSINESS_SUCCESS
		operLog.OperIp = utils.GetClientIP(c)
		operLog.OperUrl = c.Request.URI().String()
		data, ok := c.Get(api.LOGIN_INFO_KEY)
		if ok {
			if loggininfor, ok := data.(v13.LoginUser); ok {
				operLog.OperName = loggininfor.User.Username
				operLog.DeptName = loggininfor.User.Dept.DeptName
			}
		}
		// 设置方法名称
		operLog.Method = utils.B2s(c.Path())
		// 设置请求方法
		operLog.RequestMethod = utils.B2s(c.Method())
		var start, stop time.Time
		start = time.Now()
		operTime := utils.Time2Second(start)

		c.Next(ctx)

		stop = time.Now()
		// 设置消耗时间
		costTime := int64(stop.Sub(start).Seconds())
		// 设置额外信息
		buildExtraInfo(c, extraInfo, operLog)
		var response core.ErrResponse
		if err := jsoniter.Unmarshal(c.Response.Body(), &response); err != nil {
			log.Errorf("异常信息: %s", err.Error())
		}

		// 响应解析失败
		if response.Code == code.ERROR {
			operLog.Status = v12.BUSINESS_FAIL
			operLog.ErrorMsg = response.Message
		}

		operLog.OperTime = operTime
		operLog.CostTime = costTime
		if err := logSaver.SaveLog(ctx, operLog, &api.CreateOptions{}); err != nil {
			log.Errorf("异常信息: %s", err.Error())
		}
	}
}

func buildExtraInfo(c *app.RequestContext, extraInfo *OperLogExtra, operLog *v1.OperLog) {
	operLog.Title = extraInfo.Title
	operLog.BusinessType = extraInfo.BusinessType
	operLog.OperatorType = extraInfo.OperatorType
	// 需要保存请求参数
	if extraInfo.IsSaveRequestData {
		params := make(map[string]string)
		method := utils.B2s(c.Request.Method())
		visitFn := func(key []byte, value []byte) {
			k, v := utils.B2s(key), utils.B2s(value)
			if !strutil.ContainsAny(k, extraInfo.ExcludeParamNames...) {
				params[k] = v
			}
		}

		if method == http.MethodPost {
			c.Request.PostArgs().VisitAll(visitFn)
		} else {
			c.Request.URI().QueryArgs().VisitAll(visitFn)
		}

		data, _ := jsoniter.Marshal(params)
		jsonStr := utils.B2s(data)
		operLog.OperParam = jsonStr
	}

	// 需要保存响应参数
	if extraInfo.IsSaveResponseData {
		result := utils.B2s(c.Response.Body())
		operLog.JsonResult_ = result
	}
}
