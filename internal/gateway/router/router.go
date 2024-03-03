package router

import (
	"context"
	system2 "github.com/user823/Sophie/api/domain/system/v1"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/user823/Sophie/api"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/gateway/controller/v1/system"
	"github.com/user823/Sophie/internal/gateway/rpc"
	"github.com/user823/Sophie/internal/pkg/code"
	mw "github.com/user823/Sophie/internal/pkg/middleware"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/internal/pkg/middleware/secure"
	"github.com/user823/Sophie/pkg/core"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
)

type Options struct {
	Middlewares []app.HandlerFunc
	Healthz     bool
	Info        *JwtInfo
	BaseAPI     string
}

type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

type Option func(o *Options)

func WithHealthz(on bool) Option {
	return func(o *Options) {
		o.Healthz = on
	}
}

func WithMiddlewares(mws []app.HandlerFunc) Option {
	return func(o *Options) {
		o.Middlewares = mws
	}
}

func WithJwtInfo(info *JwtInfo) Option {
	return func(o *Options) {
		o.Info = info
	}
}

func WithBaseAPI(api string) Option {
	return func(o *Options) {
		o.BaseAPI = api
	}
}

func InitRouter(h *server.Hertz, opts ...Option) {
	opt := &Options{}
	for i := range opts {
		opts[i](opt)
	}

	// 获取基础组件
	logsaver := &rpcLogSaver{}
	captchaController := NewCaptchaController(false)

	// 安装通用中间件
	if h == nil {
		log.Fatalf("Insecure engine has not been prepared already: %s", system2.ServiceName)
		return
	}
	h.Use(opt.Middlewares...)

	// 未找到资源
	h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		core.WriteResponseE(c, errors.CodeMessage(code.NOT_FOUND, "未找到资源"), nil)
	})
	s := h.Group(opt.BaseAPI)

	// 认证模块
	a_ := s.Group("/auth")
	{
		jwtStrategy := newJWTAuth(opt.Info).(*auth.JWTStrategy)
		a_.POST("/login", captchaController.CheckCaptcha, jwtStrategy.LoginHandler)
		a_.DELETE("/logout", jwtStrategy.LogoutHandler)
		a_.POST("/refresh", jwtStrategy.RefreshHandler)
		a_.POST("/register", captchaController.CheckCaptcha, Register)
	}

	auto := newAutoAuth(opt.Info)
	// 系统模块
	system_ := s.Group("/system", auto.AuthFunc())
	{
		// 用户服务
		userController := system.NewUserController()
		user_ := system_.Group("/user")
		user_.GET("/list", secure.RequirePermissions("system:user:list"), userController.List)
		user_.POST("/export", secure.RequirePermissions("system:user:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), userController.Export)
		user_.POST("/importData", secure.RequirePermissions("system:user:import"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_IMPORT}), userController.ImportData)
		user_.POST("/importTemplate", userController.ImportTemplate)
		user_.GET("/getInfo", userController.GetInfo)
		user_.GET("/info/:username", userController.GetInfoWithName)
		user_.GET("/:userid", secure.RequirePermissions("system:user:query"), userController.GetInfoWithId)
		user_.GET("", secure.RequirePermissions("system:user:query"), userController.GetInfoWithId2)
		user_.POST("", secure.RequirePermissions("system:user:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), userController.Add)
		user_.PUT("", secure.RequirePermissions("system:user:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), userController.Edit)
		user_.DELETE("/:userIds", secure.RequirePermissions("system:user:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), userController.Remove)
		user_.PUT("/resetPwd", secure.RequirePermissions("system:user:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), userController.ResetPwd)
		user_.PUT("/changeStatus", secure.RequirePermissions("system:user:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), userController.ChangeStatus)
		user_.GET("/authRole/:userId", secure.RequirePermissions("system:user:query"), userController.AuthRoleWithId)
		user_.PUT("/authRole", secure.RequirePermissions("system:user:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "用户管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_GRANT}), userController.InsertAuthRole)
		user_.GET("/deptTree", secure.RequirePermissions("system:user:list"), userController.DeptTree)

		// 角色服务
		roleController := system.NewRoleController()
		role_ := system_.Group("/role")
		role_.GET("/list", secure.RequirePermissions("system:role:list"), roleController.List)
		role_.POST("/export", secure.RequirePermissions("system:role:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), roleController.Export)
		role_.GET("/:roleId", secure.RequirePermissions("system:role:query"), roleController.GetInfo)
		role_.POST("/add", secure.RequirePermissions("system:role:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), roleController.Add)
		role_.PUT("/edit", secure.RequirePermissions("system:role:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), roleController.Edit)
		role_.PUT("/dataScope", secure.RequirePermissions("system:role:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), roleController.DataScope)
		role_.PUT("/changeStatus", secure.RequirePermissions("system:role:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), roleController.ChangeStatus)
		role_.DELETE("/:roleIds", secure.RequirePermissions("system:role:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), roleController.Remove)
		role_.GET("/optionselect", secure.RequirePermissions("system:role:query"), roleController.OptionSelect)
		role_.GET("/authUser/allocatedList", secure.RequirePermissions("system:role:list"), roleController.AllocatedList)
		role_.GET("/authUser/unallocatedList", secure.RequirePermissions("system:role:list"), roleController.UnallocatedList)
		role_.PUT("/authUser/cancel", secure.RequirePermissions("system:role:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_GRANT}), roleController.CancelAuthUser)
		role_.PUT("/authUser/cancelAll", secure.RequirePermissions("system:role:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_GRANT}), roleController.CancelAuthUserAll)
		role_.PUT("/authUser/selectAll", secure.RequirePermissions("system:role:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "角色管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_GRANT}), roleController.SelectAuthUserALl)
		role_.GET("/deptTree/:roleId", secure.RequirePermissions("system:role:query"), roleController.DeptTree)

		// 个人信息服务
		profileController := system.NewProfileController()
		profile_ := user_.Group("/profile")
		profile_.GET("", profileController.Profile)
		profile_.PUT("", mw.Log(logsaver, map[string]any{mw.TITLE: "个人信息", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), profileController.UpdateProfile)
		profile_.PUT("/updatePwd", mw.Log(logsaver, map[string]any{mw.TITLE: "个人信息", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), profileController.UpdatePwd)
		profile_.POST("/avatar", mw.Log(logsaver, map[string]any{mw.TITLE: "个人信息", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), profileController.Avatar)

		// 岗位管理
		postController := system.NewPostController()
		post_ := system_.Group("/post")
		post_.GET("/list", secure.RequirePermissions("system:post:list"), postController.List)
		post_.POST("/export", secure.RequirePermissions("system:post:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "岗位管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), postController.Export)
		post_.GET("/:postId", secure.RequirePermissions("system:post:query"), postController.GetInfo)
		post_.POST("", secure.RequirePermissions("system:post:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "岗位管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), postController.Add)
		post_.PUT("", secure.RequirePermissions("system:post:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "岗位管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), postController.Edit)
		post_.DELETE("/:postIds", secure.RequirePermissions("system:post:remove"), postController.Remove)
		post_.GET("/optionselect", postController.OptionSelect)

		// 操作日志记录
		operlogController := system.NewOperlogController()
		operlog_ := system_.Group("/operlog")
		operlog_.GET("/list", secure.RequirePermissions("system:operlog:list"), operlogController.List)
		operlog_.POST("/export", secure.RequirePermissions("system:operlog:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "操作日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), operlogController.Export)
		operlog_.DELETE("/:operIds", secure.RequirePermissions("system:operlog:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "操作日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), operlogController.Remove)
		operlog_.DELETE("/clean", secure.RequirePermissions("system:operlog:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "操作日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_CLEAN}), operlogController.Clean)

		// 公告服务
		noticeController := system.NewNoticeController()
		notice_ := system_.Group("/notice")
		notice_.GET("/list", secure.RequirePermissions("system:notice:list"), noticeController.List)
		notice_.GET("/:noticeId", secure.RequirePermissions("system:notice:query"), noticeController.GetInfo)
		notice_.POST("", secure.RequirePermissions("system:notice:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "公告服务", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), noticeController.Add)
		notice_.PUT("", secure.RequirePermissions("system:notice:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "公告服务", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), noticeController.Edit)
		notice_.DELETE("/:noticeIds", secure.RequirePermissions("system:notice:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "公告服务", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), noticeController.Remove)

		// 菜单服务
		menuController := system.NewMenuController()
		menu_ := system_.Group("/menu")
		menu_.GET("/list", secure.RequirePermissions("system:menu:list"), menuController.List)
		menu_.GET("/:menuId", secure.RequirePermissions("system:menu:query"), menuController.GetInfo)
		menu_.GET("/treeselect", menuController.TreeSelect)
		menu_.GET("/roleMenuTreeselect/:roleId", menuController.RoleMenuTreeselect)
		menu_.POST("/add", secure.RequirePermissions("system:menu:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "菜单服务", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), menuController.Add)
		menu_.PUT("", secure.RequirePermissions("system:menu:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "菜单服务", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), menuController.Edit)
		menu_.DELETE("/:menuId", secure.RequirePermissions("system:menu:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "菜单服务", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), menuController.Remove)
		menu_.GET("/getRouters", menuController.GetRouters)

		// 系统访问记录
		logininfoController := system.NewLogininfoController()
		logininfo_ := system_.Group("/logininfor")
		logininfo_.GET("/list", secure.RequirePermissions("system:logininfor:list"), logininfoController.List)
		logininfo_.POST("/export", secure.RequirePermissions("system:logininfor:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "登录日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), logininfoController.Export)
		logininfo_.DELETE("/:infoIds", secure.RequirePermissions("system:logininfor:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "登录日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), logininfoController.Remove)
		logininfo_.DELETE("/clean", secure.RequirePermissions("system:logininfor:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "登录日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), logininfoController.Clean)
		logininfo_.GET("/unlock/:userName", secure.RequirePermissions("system:logininfor:unlock"), mw.Log(logsaver, map[string]any{mw.TITLE: "登录日志", mw.BUSINESSTYE: system2.BUSINESSTYPE_OTHER}), logininfoController.Unlock)

		// 字典类型管理
		dictController := system.NewDictController()
		dict_type_ := system_.Group("/dict/type")
		dict_type_.GET("/list", secure.RequirePermissions("system:dict:list"), dictController.ListType)
		dict_type_.POST("/export", secure.RequirePermissions("system:dict:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), dictController.ExportType)
		dict_type_.GET("/:dictId", secure.RequirePermissions("system:dict:query"), dictController.GetInfoType)
		dict_type_.POST("", secure.RequirePermissions("system:dict:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), dictController.AddType)
		dict_type_.PUT("", secure.RequirePermissions("system:dict:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), dictController.EditType)
		dict_type_.DELETE("/:dictIds", secure.RequirePermissions("system:dict:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), dictController.RemoveType)
		dict_type_.DELETE("/refreshCache", secure.RequirePermissions("system:dict:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_CLEAN}), dictController.RefreshCache)
		dict_type_.GET("/optionselect", dictController.OptionSelect)

		// 字典数据管理
		dict_data_ := system_.Group("/dict/data")
		dict_data_.GET("/list", secure.RequirePermissions("system:dict:list"), dictController.ListData)
		dict_data_.POST("/export", secure.RequirePermissions("system:dict:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典数据", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), dictController.ExportData)
		dict_data_.GET("/:dictCode", secure.RequirePermissions("system:dict:query"), dictController.GetInfoData)
		dict_data_.GET("/type/:dictType", dictController.DictType)
		dict_data_.POST("", secure.RequirePermissions("system:dict:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典数据", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), dictController.AddData)
		dict_data_.PUT("", secure.RequirePermissions("system:dict:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典数据", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), dictController.EditData)
		dict_data_.DELETE("/:dictCodes", secure.RequirePermissions("system:dict:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "字典数据", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), dictController.RemoveData)

		// 部门服务
		deptController := system.NewDeptController()
		dept_ := system_.Group("/dept")
		dept_.GET("/list", secure.RequirePermissions("system:dept:list"), deptController.List)
		dept_.GET("/list/exclude/:deptId", secure.RequirePermissions("system:dept:list"), deptController.ExcludeChild)
		dept_.GET("/:deptId", secure.RequirePermissions("system:dept:quert"), deptController.GetInfo)
		dept_.POST("", secure.RequirePermissions("system:dept:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "部门管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), deptController.Add)
		dept_.PUT("", secure.RequirePermissions("system:dept:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "部门管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), deptController.Edit)
		dept_.DELETE("/:deptId", secure.RequirePermissions("system:dept:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "部门管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), deptController.Remove)

		// 参数配置服务
		configController := system.NewConfigController()
		config_ := system_.Group("/config")
		config_.GET("/list", secure.RequirePermissions("system:config:list"), configController.List)
		config_.POST("/export", secure.RequirePermissions("system:config:export"), mw.Log(logsaver, map[string]any{mw.TITLE: "参数管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_EXPORT}), configController.Export)
		config_.GET("/:configId", configController.GetInfo)
		config_.GET("/configKey/:configKey", configController.GetConfigKey)
		config_.POST("", secure.RequirePermissions("system:config:add"), mw.Log(logsaver, map[string]any{mw.TITLE: "参数管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_INSERT}), configController.Add)
		config_.PUT("", secure.RequirePermissions("system:config:edit"), mw.Log(logsaver, map[string]any{mw.TITLE: "参数管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_UPDATE}), configController.Edit)
		config_.DELETE("/:configIds", secure.RequirePermissions("system:config:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "参数管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_DELETE}), configController.Remove)
		config_.DELETE("/refreshCache", secure.RequirePermissions("system:config:remove"), mw.Log(logsaver, map[string]any{mw.TITLE: "参数管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_CLEAN}), configController.RefreshCache)

		// 在线用户监控
		onlineUserController := system.NewOnlineUserController()
		onlineUser_ := system_.Group("/online")
		onlineUser_.GET("/list", secure.RequirePermissions("monitor:online:list"), onlineUserController.List)
		onlineUser_.DELETE("/:tokenId", secure.RequirePermissions("monitor:online:forceLogout"), mw.Log(logsaver, map[string]any{mw.TITLE: "参数管理", mw.BUSINESSTYE: system2.BUSINESSTYPE_FORCE}), onlineUserController.ForceLogout)
	}

	//// 文件模块
	//file_ := s.Group("/file", auto.AuthFunc())
	//{
	//	fileController := file.NewFileController()
	//	file_.POST("/upload", fileController.Upload)
	//}
	//
	//// 定时任务模块
	//schedule_ := s.Group("/schedule", auto.AuthFunc())
	//{
	//	jobController := job.NewJobController()
	//	job_ := schedule_.Group("/job")
	//	job_.GET("/list", secure.RequirePermissions("system:job:list"), jobController.List)
	//	job_.POST("/export", secure.RequirePermissions("monitor:job:export"), jobController.Export)
	//	job_.GET("/:jobId", secure.RequirePermissions("monitor:jon:query"), jobController.GetInfo)
	//	job_.POST("/", secure.RequirePermissions("monitor:job:add"), jobController.Add)
	//	job_.PUT("/", secure.RequirePermissions("monitor:job:edit"), jobController.Edit)
	//	job_.PUT("/changeStatus", secure.RequirePermissions("monitor:job:changeStatus"), jobController.ChangeStatus)
	//	job_.PUT("/run", secure.RequirePermissions("monitor:job:changeStatus"), jobController.Run)
	//	job_.DELETE("/:jobIds", secure.RequirePermissions("monitor:job:remove"), jobController.Remove)
	//}
	//
	// 代码生成模块
	code_ := s.Group("/code")
	{
		// 验证码服务
		code_.GET("", captchaController.CreateCaptcha)

		// 代码生成服务
		//genController := gen.NewGenController()
		//gen_ := code_.Group("/gen", auto.AuthFunc())
		//gen_.GET("/list", secure.RequirePermissions("tool:gen:list"), genController.List)
		//gen_.GET("/:tableId", secure.RequirePermissions("tool:gen:query"), genController.GetInfo)
		//gen_.GET("/db/list", secure.RequirePermissions("tool:gen:list"), genController.DataList)
		//gen_.GET("/column/:tableId", genController.ColumnList)
		//gen_.POST("/importTable", secure.RequirePermissions("tool:gen:import"), genController.ImportTableSave)
		//gen_.PUT("/", secure.RequirePermissions("tool:gen:edit"), genController.EditSave)
		//gen_.DELETE("/:tableIds", secure.RequirePermissions("tool:gen:remove"), genController.Remove)
		//gen_.GET("/preview/:tableId", secure.RequirePermissions("tool:gen:preview"), genController.Preview)
		//gen_.GET("/download/:tableName", secure.RequirePermissions("tool:gen:code"), genController.Download)
		//gen_.GET("/genCode/:tableName", secure.RequirePermissions("tool:gen:code"), genController.GenCode)
		//gen_.GET("/synchDb/:tableName", secure.RequirePermissions("tool:gen:edit"), genController.SynchDb)
		//gen_.GET("/batchGenCode", secure.RequirePermissions("tool:gen:code"), genController.BatchGenCode)
	}

	//健康检查
	if opt.Healthz {
		h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
			core.OK(ctx, "ok", nil)
		})
	}
}

type rpcLogSaver struct{}

func (r *rpcLogSaver) SaveLog(ctx context.Context, operLog *v1.OperLog, options *api.CreateOptions) error {
	resp, err := rpc.Remoting.CreateSysOperLog(ctx, &v1.CreateSysOperLogRequest{OperLog: operLog})
	if err != nil || resp.Code != code.SUCCESS {
		return rpc.ErrRPC
	}
	return nil
}
