package system

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	v13 "github.com/user823/Sophie/api/domain/system/v1"
	v1 "github.com/user823/Sophie/api/thrift/system/v1"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/internal/system/service"
	"github.com/user823/Sophie/internal/system/store/es"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/db/kv"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strconv"
	"strings"
)

// SystemServiceImpl implements the last service interface defined in the IDL.
type SystemServiceImpl struct{}

// ListConfigs implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListConfigs(ctx context.Context, req *v1.ListConfigsRequest) (resp *v1.ListConfigsResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), req.GetDateRange(), true)
	store, _ := es.GetESFactoryOr(nil)
	sysConfig := v1.ConfigInfo2SysConfig(req.GetConfigInfo())

	loginInfo := req.LoginUser
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewConfigs(store).SelectConfigList(ctx, sysConfig, getOpt)
	return &v1.ListConfigsResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysConfig2ConfigInfo(list.Items),
	}, nil
}

// ExportConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportConfig(ctx context.Context, req *v1.ExportConfigRequest) (resp *v1.ExportConfigResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	sysConfig := v1.ConfigInfo2SysConfig(req.GetConfigInfo())
	store, _ := es.GetESFactoryOr(nil)

	loginInfo := req.LoginUser
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewConfigs(store).SelectConfigList(ctx, sysConfig, getOpt)
	return &v1.ExportConfigResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysConfig2ConfigInfo(list.Items),
		SheetName: "参数数据",
	}, nil
}

// GetConfigById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetConfigById(ctx context.Context, id int64) (resp *v1.ConfigResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	config := service.NewConfigs(store).SelectConfigById(ctx, id, &api.GetOptions{Cache: true})
	return &v1.ConfigResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.SysConfig2ConfigInfo(config),
	}, nil
}

// GetConfigByKey implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetConfigByKey(ctx context.Context, key string) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	config := service.NewConfigs(store).SelectConfigByKey(ctx, key, &api.GetOptions{Cache: true})
	return utils.Ok(config), nil
}

// CreateConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateConfig(ctx context.Context, req *v1.CreateConfigRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	configSrv := service.NewConfigs(store)
	sysConfig := v1.ConfigInfo2SysConfig(req.GetConfigInfo())

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !configSrv.CheckConfigKeyUnique(ctx, sysConfig, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增参数 %s 失败，参数键名已存在", sysConfig.ConfigName)), nil
	}
	sysConfig.CreateBy = loginInfo.User.UserName
	if err = configSrv.InsertConfig(ctx, sysConfig, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateConfig(ctx context.Context, req *v1.UpdateConfigReqeust) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	configSrv := service.NewConfigs(store)
	sysConfig := v1.ConfigInfo2SysConfig(req.GetConfigInfo())

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !configSrv.CheckConfigKeyUnique(ctx, sysConfig, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增参数 %s 失败，参数键名已存在", sysConfig.ConfigName)), nil
	}
	sysConfig.UpdateBy = loginInfo.User.UserName
	if err = configSrv.UpdateConfig(ctx, sysConfig, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteConfig(ctx context.Context, req *v1.DeleteConfigReqeust) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if err = service.NewConfigs(store).DeleteConfigByIds(ctx, req.GetConfigIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// RefreshConfig implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RefreshConfig(ctx context.Context) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	service.NewConfigs(store).ResetConfigCache()
	return utils.Ok("操作成功"), nil
}

// ListDepts implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDepts(ctx context.Context, req *v1.ListDeptsRequest) (resp *v1.ListDeptsResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	sysDept := v1.DeptInfo2SysDept(req.GetDeptInfo())

	loginInfo := req.LoginUser
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewDepts(store).SelectDeptList(ctx, sysDept, &api.GetOptions{Cache: true})
	return &v1.ListDeptsResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MSysDept2DeptInfo(list.Items),
	}, nil
}

// ListDeptsExcludeChild implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDeptsExcludeChild(ctx context.Context, req *v1.ListDeptsExcludeChildRequest) (resp *v1.ListDeptsResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	sysDepts := service.NewDepts(store).SelectDeptList(ctx, &v13.SysDept{}, &api.GetOptions{Cache: true})
	depts := make([]*v1.DeptInfo, 0, sysDepts.TotalCount)
	loginInfo := req.LoginUser
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

c:
	for i := range sysDepts.Items {
		if sysDepts.Items[i].DeptId == req.Id {
			continue
		}
		ancestors := strings.Split(sysDepts.Items[i].Ancestors, ",")
		strid := strconv.FormatInt(req.Id, 10)
		for _, ancestor := range ancestors {
			if ancestor == strid {
				continue c
			}
		}

		depts = append(depts, v1.SysDept2DeptInfo(sysDepts.Items[i]))
	}
	return &v1.ListDeptsResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     depts,
	}, nil
}

// GetDeptById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetDeptById(ctx context.Context, req *v1.GetDeptByIdReq) (resp *v1.DeptResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	deptSrv := service.NewDepts(store)

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !deptSrv.CheckDeptDataScope(ctx, req.Id, &api.GetOptions{Cache: false}) {
		return &v1.DeptResponse{
			BaseResp: utils.Fail("没有权限访问部门数据！"),
		}, nil
	}
	return &v1.DeptResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.SysDept2DeptInfo(deptSrv.SelectDeptById(ctx, req.Id, &api.GetOptions{Cache: true})),
	}, nil
}

// CreateDept implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateDept(ctx context.Context, req *v1.CreateDeptRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	deptSrv := service.NewDepts(store)
	sysDept := v1.DeptInfo2SysDept(req.GetDept())

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !deptSrv.CheckDeptNameUnique(ctx, sysDept, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增部门 %s 失败，部门名称已存在", sysDept.DeptName)), nil
	}

	// 首先验证格式
	if err = sysDept.Validate(); err != nil {
		return utils.Fail(err.Error()), nil
	}

	sysDept.CreateBy = loginInfo.User.UserName
	if err = deptSrv.InsertDept(ctx, sysDept, &api.CreateOptions{Validate: true}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateDept implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateDept(ctx context.Context, req *v1.UpdateDeptRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	deptSrv := service.NewDepts(store)
	sysDept := v1.DeptInfo2SysDept(req.GetDept())

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !deptSrv.CheckDeptNameUnique(ctx, sysDept, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改部门 %s 失败，部门名称已存在", sysDept.DeptName)), nil
	}
	if sysDept.ParentId == sysDept.DeptId {
		return utils.Fail(fmt.Sprintf("修改部门 %s 失败，上级部门不能是自己", sysDept.DeptName)), nil
	}
	if sysDept.Status == v13.DEPTDISABLE && deptSrv.SelectNormalChildrenDeptById(ctx, sysDept.DeptId, &api.GetOptions{Cache: true}) > 0 {
		return utils.Fail("该部门包含未停用的子部门"), nil
	}

	sysDept.UpdateBy = loginInfo.User.UserName
	if err = deptSrv.UpdateDept(ctx, sysDept, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteDept implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteDept(ctx context.Context, req *v1.DeleteDeptRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	deptSrv := service.NewDepts(store)

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if deptSrv.HasChildByDeptId(ctx, req.GetDeptId(), &api.GetOptions{Cache: false}) {
		return utils.Warn("存在下级部门，无法删除"), nil
	}
	if deptSrv.CheckDeptExistUser(ctx, req.GetDeptId(), &api.GetOptions{Cache: false}) {
		return utils.Warn("部门存在用户，不允许删除"), nil
	}
	if !deptSrv.CheckDeptDataScope(ctx, req.GetDeptId(), &api.GetOptions{Cache: false}) {
		return utils.Fail("没有权限访问部门数据！"), nil
	}

	if err = deptSrv.DeleteDeptById(ctx, req.GetDeptId(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListDictDatas implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDictDatas(ctx context.Context, req *v1.ListDictDatasRequest) (resp *v1.ListDictDatasResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	sysDictData := v1.DictData2SysDictData(req.GetDictData())

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewDictDatas(store).SelectDictDataList(ctx, sysDictData, getOpt)
	return &v1.ListDictDatasResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysDictData2DictData(list.Items),
	}, nil
}

// ExportDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportDictData(ctx context.Context, req *v1.ExportDictDataRequest) (resp *v1.ExportDictDataResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	sysDictData := v1.DictData2SysDictData(req.GetDictData())

	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	list := service.NewDictDatas(store).SelectDictDataList(ctx, sysDictData, &api.GetOptions{Cache: true})
	return &v1.ExportDictDataResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysDictData2DictData(list.Items),
		SheetName: "字典数据",
	}, nil
}

// GetDictDataByCode implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetDictDataByCode(ctx context.Context, code int64) (resp *v1.DictDataResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	sysDictData := service.NewDictDatas(store).SelectDictDataById(ctx, code, &api.GetOptions{Cache: true})
	return &v1.DictDataResponse{
		BaseResp: utils.Ok("操作成功"),
		DictData: v1.SysDictData2DictData(sysDictData),
	}, nil
}

// ListDictDataByType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDictDataByType(ctx context.Context, dictType string) (resp *v1.ListDictDatasResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	list := service.NewDictTypes(store).SelectDictDataByType(ctx, dictType, &api.GetOptions{Cache: true})
	return &v1.ListDictDatasResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysDictData2DictData(list.Items),
	}, nil
}

// CreateDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateDictData(ctx context.Context, req *v1.CreateDictDataRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysDictData := v1.DictData2SysDictData(req.GetDictData())

	// 首先验证格式
	if err = sysDictData.Validate(); err != nil {
		return utils.Fail(err.Error()), nil
	}

	sysDictData.CreateBy = loginInfo.User.UserName
	if err = service.NewDictDatas(store).InsertDictData(ctx, sysDictData, &api.CreateOptions{Validate: true}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateDictData(ctx context.Context, req *v1.UpdateDictDataRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysDictData := v1.DictData2SysDictData(req.GetDictData())
	sysDictData.UpdateBy = loginInfo.User.UserName
	if err = service.NewDictDatas(store).UpdateDictData(ctx, sysDictData, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteDictData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteDictData(ctx context.Context, req *v1.DeleteDictDataRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewDictDatas(store).DeleteDictDataByIds(ctx, req.GetDictCodes(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListDictTypes implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDictTypes(ctx context.Context, req *v1.ListDictTypesRequest) (resp *v1.ListDictTypesResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	getOpt := utils.BuildGetOption(req.GetPageInfo(), req.GetDateRange(), true)
	sysDictType := v1.DictType2SysDictType(req.GetDictType())
	list := service.NewDictTypes(store).SelectDictTypeList(ctx, sysDictType, getOpt)
	return &v1.ListDictTypesResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysDictType2DictType(list.Items),
	}, nil
}

// ExportDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportDictType(ctx context.Context, req *v1.ExportDictTypeRequest) (resp *v1.ExportDictTypeResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysDictType := v1.DictType2SysDictType(req.GetDictType())
	list := service.NewDictTypes(store).SelectDictTypeList(ctx, sysDictType, &api.GetOptions{Cache: true})
	return &v1.ExportDictTypeResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysDictType2DictType(list.Items),
		SheetName: "字典类型",
	}, nil
}

// GetDictTypeById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetDictTypeById(ctx context.Context, id int64) (resp *v1.DictTypeResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	sysDictType := service.NewDictTypes(store).SelectDictTypeById(ctx, id, &api.GetOptions{Cache: true})
	return &v1.DictTypeResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.SysDictType2DictType(sysDictType),
	}, nil
}

// CreateDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateDictType(ctx context.Context, req *v1.CreateDictTypeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysDictType := v1.DictType2SysDictType(req.GetDictType())
	dictTypeSrv := service.NewDictTypes(store)
	if !dictTypeSrv.CheckDictTypeUnique(ctx, sysDictType, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增字典 %s 失败，字典类型已存在", sysDictType.DictName)), nil
	}

	// 验证格式
	if err = sysDictType.Validate(); err != nil {
		return utils.Fail(err.Error()), nil
	}

	sysDictType.CreateBy = loginInfo.User.UserName
	if err = dictTypeSrv.InsertDictType(ctx, sysDictType, &api.CreateOptions{Validate: true}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateDictType(ctx context.Context, req *v1.UpdateDictTypeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	dictTypeSrv := service.NewDictTypes(store)
	sysDictType := v1.DictType2SysDictType(req.GetDictType())
	if !dictTypeSrv.CheckDictTypeUnique(ctx, sysDictType, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改字典 %s 失败，字典类型已存在", sysDictType.DictName)), nil
	}

	sysDictType.UpdateBy = loginInfo.User.UserName
	if err = dictTypeSrv.UpdateDictType(ctx, sysDictType, &api.UpdateOptions{}); err != nil {
		log.Infof("修改字典类型错误: %s", err.Error())
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteDictType(ctx context.Context, req *v1.DeleteDictTypeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewDictTypes(store).DeleteDictTypeByIds(ctx, req.GetDictIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// RefreshDictType implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RefreshDictType(ctx context.Context) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	service.NewDictTypes(store).ResetDictCache(ctx)
	return utils.Ok("操作成功"), nil
}

// DictTypeOptionSelect implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DictTypeOptionSelect(ctx context.Context) (resp *v1.DictTypeOptionSelectResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	list := service.NewDictTypes(store).SelectDictTypeAll(ctx, &api.GetOptions{Cache: false})
	return &v1.DictTypeOptionSelectResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MSysDictType2DictType(list.Items),
	}, nil
}

// ListSysLogininfos implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysLogininfos(ctx context.Context, req *v1.ListSysLogininfosRequest) (resp *v1.ListSysLogininfosResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	getOpt := utils.BuildGetOption(req.GetPageInfo(), req.GetDateRange(), true)
	sysLogininfor := v1.LoginInfo2SysLogininfo(req.GetLoginInfo())

	list := service.NewLogininfors(store).SelectLogininforList(ctx, sysLogininfor, getOpt)
	return &v1.ListSysLogininfosResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysLogininfor2Logininfor(list.Items),
	}, nil
}

// ExportLogininfo implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportLogininfo(ctx context.Context, req *v1.ExportLogininfoRequest) (resp *v1.ExportLogininfoResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysLogininfor := v1.LoginInfo2SysLogininfo(req.GetLoginInfo())

	list := service.NewLogininfors(store).SelectLogininforList(ctx, sysLogininfor, &api.GetOptions{Cache: true})
	return &v1.ExportLogininfoResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysLogininfor2Logininfor(list.Items),
		SheetName: "登录日志",
	}, nil
}

// RemoveSysLogininfosById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RemoveSysLogininfosById(ctx context.Context, req *v1.RemoveSysLogininfosByIdRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewLogininfors(store).DeleteLogininforByIds(ctx, req.GetInfoIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// LogininfoClean implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) LogininfoClean(ctx context.Context) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	if err = service.NewLogininfors(store).CleanLogininfor(ctx, &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UnlockByUserName implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UnlockByUserName(ctx context.Context, username string) (resp *v1.BaseResp, err error) {
	redisCli := kv.NewKVStore("redis", nil).(kv.RedisStore)
	redisCli.SetKeyPrefix(kv.PWD_ERR_CNT_KEY)
	redisCli.DeleteKey(ctx, username)
	return utils.Ok("操作成功"), nil
}

// CreateSysLogininfo implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysLogininfo(ctx context.Context, req *v1.CreateSysLogininfoRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysLogininfor := v1.LoginInfo2SysLogininfo(req.GetLoginInfo())
	if err = service.NewLogininfors(store).InsertLogininfor(ctx, sysLogininfor, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListSysMenus implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysMenus(ctx context.Context, req *v1.ListSysMenusRequest) (resp *v1.ListSysMenusResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysMenu := v1.MenuInfo2SysMenu(req.GetMenuInfo())

	menus := service.NewMenus(store).SelectMenuListWithMenu(ctx, sysMenu, loginInfo.User.UserId, &api.GetOptions{Cache: true})
	return &v1.ListSysMenusResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MSysMenu2MenuInfo(menus.Items),
	}, nil
}

// GetSysMenuById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysMenuById(ctx context.Context, id int64) (resp *v1.SysMenuResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	menu := service.NewMenus(store).SelectMenuById(ctx, id, &api.GetOptions{Cache: true})
	return &v1.SysMenuResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.SysMenu2MenuInfo(menu),
	}, nil
}

// ListTreeMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListTreeMenu(ctx context.Context, req *v1.ListTreeMenuRequest) (resp *v1.ListTreeMenuResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysMenu := v1.MenuInfo2SysMenu(req.GetMenuInfo())
	menuSrv := service.NewMenus(store)

	menus := menuSrv.SelectMenuListWithMenu(ctx, sysMenu, loginInfo.User.UserId, &api.GetOptions{Cache: true})
	res := menuSrv.BuildMenuTreeSelect(ctx, menus.Items)

	return &v1.ListTreeMenuResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MTreeSelectTrans(res),
	}, nil
}

// ListTreeMenuByRoleid implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListTreeMenuByRoleid(ctx context.Context, req *v1.ListTreeMenuByRoleidRequest) (resp *v1.RoleMenuResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	menuSrv := service.NewMenus(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := menuSrv.SelectMenuList(ctx, loginInfo.User.UserId, &api.GetOptions{Cache: true})
	keys := menuSrv.SelectMenuListByRoleId(ctx, req.Id, &api.GetOptions{Cache: true})
	menus := menuSrv.BuildMenuTreeSelect(ctx, list.Items)
	return &v1.RoleMenuResponse{
		BaseResp:    utils.Ok("操作成功"),
		CheckedKeys: keys,
		Menus:       v1.MTreeSelectTrans(menus),
	}, nil
}

// CreateMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateMenu(ctx context.Context, req *v1.CreateMenuRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	menuSrv := service.NewMenus(store)
	sysMenu := v1.MenuInfo2SysMenu(req.GetMenuInfo())
	if !menuSrv.CheckMenuNameUnique(ctx, sysMenu, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增菜单 %s 失败，菜单名称已存在", sysMenu.MenuName)), nil
	}
	if sysMenu.IsFrame == v13.YES_FRAME && !strutil.IsHttp(sysMenu.Path) {
		return utils.Fail(fmt.Sprintf("新增菜单 %s 失败，地址必须以http(s)://开头", sysMenu.MenuName)), nil
	}

	sysMenu.CreateBy = loginInfo.User.UserName
	if err = menuSrv.InsertMenu(ctx, sysMenu, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateMenu(ctx context.Context, req *v1.UpdateMenuRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	menuSrv := service.NewMenus(store)
	sysMenu := v1.MenuInfo2SysMenu(req.GetMenuInfo())
	if !menuSrv.CheckMenuNameUnique(ctx, sysMenu, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增菜单 %s 失败，菜单名称已存在", sysMenu.MenuName)), nil
	}
	if sysMenu.IsFrame == v13.YES_FRAME && !strutil.IsHttp(sysMenu.Path) {
		return utils.Fail(fmt.Sprintf("新增菜单 %s 失败，地址必须以http(s)://开头", sysMenu.MenuName)), nil
	}
	if sysMenu.MenuId == sysMenu.ParentId {
		return utils.Fail(fmt.Sprintf("修改菜单 %s 失败，上级菜单不能选择自己", sysMenu.MenuName)), nil
	}

	sysMenu.UpdateBy = loginInfo.User.UserName
	if err = menuSrv.UpdateMenu(ctx, sysMenu, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteMenu implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteMenu(ctx context.Context, req *v1.DeleteMenuRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	menuSrv := service.NewMenus(store)
	if menuSrv.HasChildByMenuId(ctx, req.GetMenuId(), &api.GetOptions{Cache: false}) {
		return utils.Fail("存在子菜单，不允许删除"), nil
	}
	if menuSrv.CheckMenuExistsRole(ctx, req.GetMenuId(), &api.GetOptions{Cache: false}) {
		return utils.Warn("菜单已分配，不允许删除"), nil
	}
	if err = menuSrv.DeleteMenuBuId(ctx, req.GetMenuId(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// GetRouters implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetRouters(ctx context.Context, req *v1.GetRoutersRequest) (resp *v1.RoutersResonse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	menuSrv := service.NewMenus(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	menus := menuSrv.SelectMenuTreeByUserId(ctx, loginInfo.User.UserId, &api.GetOptions{Cache: true})
	routers := menuSrv.BuildMenus(ctx, menus.Items)
	return &v1.RoutersResonse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MRouterVoTrans(routers),
	}, nil
}

// ListSysNotices implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysNotices(ctx context.Context, req *v1.ListSysNoticesRequest) (resp *v1.ListSysNoticesResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	sysNotice := v1.NoticeInfo2SysNotice(req.GetNoticeInfo())
	list := service.NewNotices(store).SelectNoticeList(ctx, sysNotice, getOpt)
	return &v1.ListSysNoticesResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysNotice2NoticeInfo(list.Items),
	}, nil
}

// GetSysNoticeById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysNoticeById(ctx context.Context, id int64) (resp *v1.SysNoticeResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	sysNotice := service.NewNotices(store).SelectNoticeById(ctx, id, &api.GetOptions{Cache: true})
	return &v1.SysNoticeResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.SysNotice2NoticeInfo(sysNotice),
	}, nil
}

// CreateSysNotice implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysNotice(ctx context.Context, req *v1.CreateSysNoticeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysNotice := v1.NoticeInfo2SysNotice(req.GetNoticeInfo())
	sysNotice.CreateBy = loginInfo.User.UserName
	if err = service.NewNotices(store).InsertNotice(ctx, sysNotice, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteSysNotice implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysNotice(ctx context.Context, req *v1.DeleteSysNoticeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewNotices(store).DeleteNoticeByIds(ctx, req.GetNoticeIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateSysNotice implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysNotice(ctx context.Context, req *v1.UpdateSysNoticeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysNotice := v1.NoticeInfo2SysNotice(req.GetNoticeInfo())

	sysNotice.UpdateBy = loginInfo.User.UserName
	if err = service.NewNotices(store).UpdateNotice(ctx, sysNotice, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListSysOperLogs implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysOperLogs(ctx context.Context, req *v1.ListSysOperLogsRequest) (resp *v1.ListSysOperLogsResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysOperLog := v1.OperLog2SysOperLog(req.GetOperLog())

	list := service.NewOperLogs(store).SelectOperLogList(ctx, sysOperLog, getOpt)
	return &v1.ListSysOperLogsResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysOperLog2Operlog(list.Items),
	}, nil
}

// ExportSysOperLog implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysOperLog(ctx context.Context, req *v1.ExportSysOperLogRequest) (resp *v1.ExportSysOperLogResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysOperLog := v1.OperLog2SysOperLog(req.GetOperLog())

	list := service.NewOperLogs(store).SelectOperLogList(ctx, sysOperLog, &api.GetOptions{Cache: true})
	return &v1.ExportSysOperLogResponse{
		BaseResp:  utils.Ok("操作成功"),
		OperLogs:  v1.MSysOperLog2Operlog(list.Items),
		SheetName: "操作日志",
	}, nil
}

// DeleteSysOperLog implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysOperLog(ctx context.Context, req *v1.DeleteSysOperLogRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewOperLogs(store).DeleteOperLogByIds(ctx, req.GetOperIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// OperLogClean implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) OperLogClean(ctx context.Context) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	if err = service.NewOperLogs(store).CleanOperLog(ctx, &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// CreateSysOperLog implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysOperLog(ctx context.Context, req *v1.CreateSysOperLogRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysOperLog := v1.OperLog2SysOperLog(req.GetOperLog())
	if err = service.NewOperLogs(store).InsertOperLog(ctx, sysOperLog, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListSysPosts implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysPosts(ctx context.Context, req *v1.ListSysPostsRequest) (resp *v1.ListSysPostsResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysPost := v1.PostInfo2SysPost(req.GetPostInfo())

	list := service.NewPosts(store).SelectPostList(ctx, sysPost, getOpt)
	return &v1.ListSysPostsResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysPost2PostInfo(list.Items),
	}, nil
}

// ExportSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysPost(ctx context.Context, req *v1.ExportSysPostRequest) (resp *v1.ExportSysPostResponse, err error) {
	getOpts := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	store, _ := es.GetESFactoryOr(nil)
	sysPost := v1.PostInfo2SysPost(req.PostInfo)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewPosts(store).SelectPostList(ctx, sysPost, getOpts)
	return &v1.ExportSysPostResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysPost2PostInfo(list.Items),
		SheetName: "岗位数据",
	}, nil
}

// GetSysPostById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysPostById(ctx context.Context, id int64) (resp *v1.SysPostResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	post := service.NewPosts(store).SelectPostById(ctx, id, &api.GetOptions{Cache: true})
	return &v1.SysPostResponse{
		BaseResp: utils.Ok("操作成功"),
		PostInfo: v1.SysPost2PostInfo(post),
	}, nil
}

// CreateSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysPost(ctx context.Context, req *v1.CreateSysPostRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	postSrv := service.NewPosts(store)
	sysPost := v1.PostInfo2SysPost(req.GetPostInfo())

	if !postSrv.CheckPostNameUnique(ctx, sysPost, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增岗位 %s 失败，岗位名称已存在", sysPost.PostName)), nil
	}
	if !postSrv.CheckPostCodeUnique(ctx, sysPost, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增岗位 %s 失败，岗位编码已存在", sysPost.PostCode)), nil
	}
	sysPost.CreateBy = loginInfo.User.UserName
	if err = postSrv.InsertPost(ctx, sysPost, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysPost(ctx context.Context, req *v1.UpdateSysPostRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	postSrv := service.NewPosts(store)
	sysPost := v1.PostInfo2SysPost(req.GetPostInfo())

	if !postSrv.CheckPostNameUnique(ctx, sysPost, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增岗位 %s 失败，岗位名称已存在", sysPost.PostName)), nil
	}
	if !postSrv.CheckPostCodeUnique(ctx, sysPost, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增岗位 %s 失败，岗位编码已存在", sysPost.PostCode)), nil
	}
	sysPost.UpdateBy = loginInfo.User.UserName
	if err = postSrv.UpdatePost(ctx, sysPost, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteSysPost implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysPost(ctx context.Context, req *v1.DeleteSysPostRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	postSrv := service.NewPosts(store)

	for i := range req.PostIds {
		post := postSrv.SelectPostById(ctx, req.PostIds[i], &api.GetOptions{Cache: false})
		if postSrv.CountUserPostById(ctx, req.PostIds[i], &api.GetOptions{Cache: true}) > 0 {
			return utils.Fail(fmt.Sprintf("%s 岗位已经分配, 不能删除", post.PostName)), nil
		}
	}

	if err = service.NewPosts(store).DeletePostByIds(ctx, req.GetPostIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// PostOptionSelect implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) PostOptionSelect(ctx context.Context) (resp *v1.PostOptionSelectResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	list := service.NewPosts(store).SelectPostAll(ctx, &api.GetOptions{Cache: true})
	return &v1.PostOptionSelectResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MSysPost2PostInfo(list.Items),
	}, nil
}

// Profile implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) Profile(ctx context.Context, req *v1.ProfileRequest) (resp *v1.ProfileResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	sysUser := userSrv.SelectUserByUserName(ctx, &v13.SysUser{Username: loginInfo.User.UserName}, &api.GetOptions{Cache: true})
	roleGroup := userSrv.SelectUserRoleGroup(ctx, loginInfo.User.UserName, &api.GetOptions{Cache: true})
	postGroup := userSrv.SelectUserPostGroup(ctx, loginInfo.User.UserName, &api.GetOptions{Cache: true})
	return &v1.ProfileResponse{
		BaseResp:  utils.Ok("操作成功"),
		UserInfo:  v1.SysUser2UserInfo(sysUser),
		RoleGroup: roleGroup,
		PostGroup: postGroup,
	}, nil
}

// UpdateProfile implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	userSrv := service.NewUsers(store)

	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	if sysUser.Phonenumber != "" && !userSrv.CheckPhoneUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改用户 %s 失败，手机号已存在", sysUser.Username)), nil
	}
	if sysUser.Email != "" && !userSrv.CheckEmailUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改用户 %s 失败，邮箱账号已存在", sysUser.Username)), nil
	}

	if err = userSrv.UpdateUsrProfile(ctx, sysUser, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdatePassword implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdatePassword(ctx context.Context, req *v1.UpdatePasswordRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	sysUser := userSrv.SelectUserByUserName(ctx, &v13.SysUser{Username: loginInfo.User.UserName}, &api.GetOptions{Cache: true})
	if err = auth.Compare(sysUser.Password, req.OldPassword); err != nil {
		return utils.Fail("修改密码失败，旧密码错误"), nil
	}
	if req.NewPassword_ == req.OldPassword {
		return utils.Fail("新密码不能和旧密码相同"), nil
	}
	if err = userSrv.ResetUserPwd(ctx, sysUser.Username, req.NewPassword_, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysRole(ctx context.Context, req *v1.ListSysRolesRequest) (resp *v1.ListSysRolesResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), req.DateRange, true)
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysRole := v1.RoleInfo2SysRole(req.GetRoleInfo())

	list := service.NewRoles(store).SelectRoleList(ctx, sysRole, getOpt)
	return &v1.ListSysRolesResponse{
		BaseResp: utils.Ok("操作成功"),
		Rows:     v1.MSysRole2RoleInfo(list.Items),
		Total:    list.TotalCount,
	}, nil
}

// ExportSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysRole(ctx context.Context, req *v1.ExportSysRoleRequest) (resp *v1.ExportSysRoleResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysRole := v1.RoleInfo2SysRole(req.GetRoleInfo())

	list := service.NewRoles(store).SelectRoleList(ctx, sysRole, &api.GetOptions{Cache: true})
	return &v1.ExportSysRoleResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysRole2RoleInfo(list.Items),
		SheetName: "角色数据",
	}, nil
}

// GetSysRoleByid implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetSysRoleByid(ctx context.Context, id int64) (resp *v1.SysRoleResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	role := service.NewRoles(store).SelectRoleById(ctx, id, &api.GetOptions{Cache: true})
	return &v1.SysRoleResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.SysRole2RoleInfo(role),
	}, nil
}

// CreateSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysRole(ctx context.Context, req *v1.CreateSysRoleRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	roleSrv := service.NewRoles(store)
	sysRole := v1.RoleInfo2SysRole(req.GetRoleInfo())
	if !roleSrv.CheckRoleNameUnique(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增角色 %s 失败，角色名称已存在y", sysRole.RoleName)), nil
	}
	if !roleSrv.CheckRoleKeyUnique(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增角色 %s 失败，角色权限已存在", sysRole.RoleName)), nil
	}

	sysRole.CreateBy = loginInfo.User.UserName
	if err = roleSrv.InsertRole(ctx, sysRole, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), err
	}
	return utils.Ok("操作成功"), nil
}

// UpdateSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysRole(ctx context.Context, req *v1.UpdateSysRoleRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	roleSrv := service.NewRoles(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysRole := v1.RoleInfo2SysRole(req.GetRoleInfo())

	if !roleSrv.CheckRoleAllowed(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改角色 %s 失败，角色名称已存在", sysRole.RoleName)), nil
	}
	if !roleSrv.CheckRoleKeyUnique(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增角色 %s 失败，角色权限已存在", sysRole.RoleName)), nil
	}
	sysRole.UpdateBy = loginInfo.User.UserName
	if err = roleSrv.UpdateRole(ctx, sysRole, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DataScope implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DataScope(ctx context.Context, req *v1.DataScopeRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	roleSrv := service.NewRoles(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysRole := v1.RoleInfo2SysRole(req.GetRoleInfo())

	if !roleSrv.CheckRoleAllowed(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改角色 %s 失败，角色名称已存在", sysRole.RoleName)), nil
	}
	if !roleSrv.CheckRoleKeyUnique(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增角色 %s 失败，角色权限已存在", sysRole.RoleName)), nil
	}

	if err = roleSrv.AuthDataScope(ctx, sysRole, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ChangeSysRoleStatus implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ChangeSysRoleStatus(ctx context.Context, req *v1.ChangeSysRoleStatusRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	roleSrv := service.NewRoles(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysRole := v1.RoleInfo2SysRole(req.GetRoleInfo())

	if !roleSrv.CheckRoleAllowed(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("修改角色 %s 失败，角色名称已存在", sysRole.RoleName)), nil
	}
	if !roleSrv.CheckRoleKeyUnique(ctx, sysRole, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增角色 %s 失败，角色权限已存在", sysRole.RoleName)), nil
	}

	sysRole.UpdateBy = loginInfo.User.UserName
	if err = roleSrv.UpdateRoleStatus(ctx, sysRole, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeleteSysRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysRole(ctx context.Context, req *v1.DeleteSysRoleRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewRoles(store).DeleteRoleByIds(ctx, req.GetRoleIds(), &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ListRoleOption implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListRoleOption(ctx context.Context) (resp *v1.ListSysRolesResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	list := service.NewRoles(store).SelectRoleAll(ctx, &api.GetOptions{Cache: true})
	return &v1.ListSysRolesResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysRole2RoleInfo(list.Items),
	}, nil
}

// AllocatedList implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) AllocatedList(ctx context.Context, req *v1.AllocatedListRequest) (resp *v1.ListSysUsersResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())

	list := service.NewUsers(store).SelectAllocatedList(ctx, sysUser, getOpt)
	return &v1.ListSysUsersResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysUser2UserInfo(list.Items),
	}, nil
}

// UnallocatedList implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UnallocatedList(ctx context.Context, req *v1.UnallocatedListRequest) (resp *v1.ListSysUsersResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())

	list := service.NewUsers(store).SelectUnallocatedList(ctx, sysUser, getOpt)
	return &v1.ListSysUsersResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysUser2UserInfo(list.Items),
	}, nil
}

// CancelAuthUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CancelAuthUser(ctx context.Context, req *v1.CancelAuthUserRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewRoles(store).DeleteAuthUser(ctx, &v13.SysUserRole{UserId: req.GetUserId(), RoleId: req.GetRoleId()}, &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// CancelAuthUserAll implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CancelAuthUserAll(ctx context.Context, req *v1.CancelAuthUserAllRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	if err = service.NewRoles(store).DeleteAuthUsers(ctx, req.RoleId, req.UserIds, &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// SelectAuthUserAll implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) SelectAuthUserAll(ctx context.Context, req *v1.SelectAuthUserAllRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	roleSrv := service.NewRoles(store)
	if !roleSrv.CheckRoleDataScope(ctx, req.RoleId, &api.GetOptions{Cache: false}) {
		return utils.Fail("没有权限访问角色数据！"), nil
	}
	if err = roleSrv.InsertAuthUsers(ctx, req.GetRoleId(), req.GetUserIds(), &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// DeptTreeByRoleId implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeptTreeByRoleId(ctx context.Context, req *v1.DeptTreeByRoleIdRequest) (resp *v1.DeptTreeByRoleIdResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	deptSrv := service.NewDepts(store)
	keys := deptSrv.SelectDeptListByRoleId(ctx, req.Id, &api.GetOptions{Cache: true})
	depts := deptSrv.SelectDeptTreeList(ctx, &v13.SysDept{}, &api.GetOptions{Cache: true})
	return &v1.DeptTreeByRoleIdResponse{
		BaseResp:    utils.Ok("操作成功"),
		CheckedKeys: keys,
		Depts:       v1.MTreeSelectTrans(depts),
	}, nil
}

// ListSysUsers implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysUsers(ctx context.Context, req *v1.ListSysUsersRequest) (resp *v1.ListSysUsersResponse, err error) {
	getOpt := utils.BuildGetOption(req.PageInfo, nil, true)
	store, _ := es.GetESFactoryOr(nil)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewUsers(store).SelectUserList(ctx, sysUser, getOpt)

	return &v1.ListSysUsersResponse{
		BaseResp: utils.Ok("操作成功"),
		Total:    list.TotalCount,
		Rows:     v1.MSysUser2UserInfo(list.Items),
	}, nil
}

// ExportSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ExportSysUser(ctx context.Context, req *v1.ExportSysUserRequest) (resp *v1.ExportSysUserResponse, err error) {
	getOpt := utils.BuildGetOption(req.GetPageInfo(), nil, true)
	store, _ := es.GetESFactoryOr(nil)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	list := service.NewUsers(store).SelectUserList(ctx, sysUser, getOpt)

	return &v1.ExportSysUserResponse{
		BaseResp:  utils.Ok("操作成功"),
		List:      v1.MSysUser2UserInfo(list.Items),
		SheetName: "用户数据",
	}, nil
}

// ImportUserData implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ImportUserData(ctx context.Context, req *v1.ImportUserDataRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	users := v1.MUserInfo2SysUser(req.GetUsers())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	msg, err := service.NewUsers(store).ImportUser(ctx, users, req.IsUpdateSupport, req.OperName, &api.UpdateOptions{})
	if err != nil {
		return utils.Fail(err.Error()), nil
	}
	return utils.Ok(msg), nil
}

// GetUserInfo implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetUserInfo(ctx context.Context, id int64) (resp *v1.UserInfoResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)

	user := service.NewUsers(store).SelectUserById(ctx, id, &api.GetOptions{Cache: true})
	if user == nil {
		return &v1.UserInfoResponse{
			BaseResp: utils.Fail("未找到用户"),
		}, nil
	}
	permissionSrv := service.NewPermissions(store)
	roles := permissionSrv.GetRolePermission(ctx, user, &api.GetOptions{Cache: true})
	permissions := permissionSrv.GetMenuPermission(ctx, user, &api.GetOptions{Cache: true})
	return &v1.UserInfoResponse{
		BaseResp:    utils.Ok("操作成功"),
		Roles:       roles,
		Permissions: permissions,
		Data:        v1.SysUser2UserInfo(user),
	}, nil
}

// GetUserInfoByName implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetUserInfoByName(ctx context.Context, name string) (resp *v1.UserInfoResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)

	user := service.NewUsers(store).SelectUserByUserName(ctx, &v13.SysUser{Username: name}, &api.GetOptions{Cache: true})
	if user == nil {
		return &v1.UserInfoResponse{
			BaseResp: utils.Fail("用户名或密码错误"),
			Data:     nil,
		}, nil
	}
	// 查询角色集合
	permissionSrv := service.NewPermissions(store)
	roles := permissionSrv.GetRolePermission(ctx, user, &api.GetOptions{Cache: true})
	permissions := permissionSrv.GetMenuPermission(ctx, user, &api.GetOptions{Cache: true})
	return &v1.UserInfoResponse{
		BaseResp:    utils.Ok("操作成功"),
		Data:        v1.SysUser2UserInfo(user),
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

// RegisterSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) RegisterSysUser(ctx context.Context, req *v1.RegisterSysUserRequest) (resp *v1.RegisterSysUserResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	// 不走es缓存
	user := v1.UserInfo2SysUser(req.GetUserInfo())
	if service.NewConfigs(store).SelectConfigByKey(ctx, "sys.account.registerUser", &api.GetOptions{Cache: true}) != "true" {
		return &v1.RegisterSysUserResponse{
			BaseResp: utils.Fail("当前系统没有开启注册功能！"),
			IsOk:     false,
		}, nil
	}
	if !service.NewUsers(store).CheckUserNameUnique(ctx, user, &api.GetOptions{Cache: false}) {
		return &v1.RegisterSysUserResponse{
			BaseResp: utils.Fail(fmt.Sprintf("保存用户 %s 失败，注册账号已经存在", user.Username)),
			IsOk:     false,
		}, nil
	}
	return &v1.RegisterSysUserResponse{
		BaseResp: utils.Ok("操作成功"),
		IsOk:     true,
	}, nil
}

// GetUserInfoById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetUserInfoById(ctx context.Context, req *v1.GetUserInfoByIdRequest) (resp *v1.UserInfoByIdResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	postSrv := service.NewPosts(store)
	roleSrv := service.NewRoles(store)
	loginuser := req.GetUser()
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginuser)

	roleList := roleSrv.SelectRoleAll(ctx, &api.GetOptions{Cache: true})
	// 搜索所有角色
	roles := roleList.Items
	// 搜索所有岗位
	// 搜索所有岗位
	posts := postSrv.SelectPostAll(ctx, &api.GetOptions{Cache: true})
	if req.Id != -1 {
		// 查询用户权限
		if !userSrv.CheckUserDataScope(ctx, req.Id, &api.GetOptions{Cache: false}) {
			return &v1.UserInfoByIdResponse{
				BaseResp: utils.Fail("没有权限访问用户数据！"),
			}, nil
		}

		// 过滤角色
		if !v13.IsUserAdmin(req.Id) {
			j := 0
			for i := 0; i < len(roleList.Items); i++ {
				if !roles[i].IsAdmin() {
					roles[j] = roles[i]
					j++
				}
			}
			roles = roles[:j]
		}

		// 该用户对应岗位和角色
		sysUser := userSrv.SelectUserById(ctx, req.Id, &api.GetOptions{Cache: true})
		postIds := postSrv.SelectPostListByUserId(ctx, req.Id, &api.GetOptions{Cache: true})
		var roleIds []int64
		for i := range sysUser.Roles {
			roleIds = append(roleIds, sysUser.Roles[i].RoleId)
		}
		return &v1.UserInfoByIdResponse{
			BaseResp: utils.Ok("操作成功"),
			Data:     v1.SysUser2UserInfo(sysUser),
			PostIds:  postIds,
			RoleIds:  roleIds,
			Roles:    v1.MSysRole2RoleInfo(roles),
			Posts:    v1.MSysPost2PostInfo(posts.Items),
		}, nil
	}
	return &v1.UserInfoByIdResponse{
		BaseResp: utils.Ok("操作成功"),
		Roles:    v1.MSysRole2RoleInfo(roles),
		Posts:    v1.MSysPost2PostInfo(posts.Items),
	}, nil
}

// CreateSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) CreateSysUser(ctx context.Context, req *v1.CreateSysUserRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if err = sysUser.Validate(); err != nil {
		return utils.Fail(err.Error()), nil
	}

	if !userSrv.CheckUserNameUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增用户 %s 失败，登录账号已存在", sysUser.Username)), nil
	}
	if sysUser.Phonenumber != "" && !userSrv.CheckPhoneUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增用户 %s 失败，手机号码已存在", sysUser.Username)), nil
	}
	if sysUser.Email != "" && !userSrv.CheckEmailUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增用户 %s 失败，邮箱账号已存在", sysUser.Username)), nil
	}

	sysUser.CreateBy = loginInfo.User.UserName
	encryptedPasswd, _ := auth.Encrypt(sysUser.Password)
	sysUser.Password = encryptedPasswd
	if err = userSrv.InsertUser(ctx, sysUser, &api.CreateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), err
	}
	return utils.Ok("操作成功"), nil
}

// UpdateSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateSysUser(ctx context.Context, req *v1.UpdateSysUserRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !userSrv.CheckUserAllowed(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail("不允许操作超级管理员用户"), nil
	}
	if !userSrv.CheckUserDataScope(ctx, sysUser.UserId, &api.GetOptions{Cache: false}) {
		return utils.Fail("没有权限访问用户数据！"), nil
	}
	if !userSrv.CheckUserNameUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增用户 %s 失败，登录账号已存在", sysUser.Username)), nil
	}
	if sysUser.Phonenumber != "" && !userSrv.CheckPhoneUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增用户 %s 失败，手机号码已存在", sysUser.Username)), nil
	}
	if sysUser.Email != "" && !userSrv.CheckEmailUnique(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail(fmt.Sprintf("新增用户 %s 失败，邮箱账号已存在", sysUser.Username)), nil
	}

	sysUser.UpdateBy = loginInfo.User.UserName
	if err = userSrv.UpdateUser(ctx, sysUser, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), err
	}
	return utils.Ok("操作成功"), nil
}

// DeleteSysUser implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) DeleteSysUser(ctx context.Context, req *v1.DeleteSysUserRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	// 排除掉当前用户
	for i := 0; i < len(req.UserIds); i++ {
		if req.UserIds[i] == loginInfo.User.UserId {
			req.UserIds = append(req.UserIds[:i], req.UserIds[i+1:]...)
			break
		}
	}

	if err = userSrv.DeleteUserByIds(ctx, req.UserIds, &api.DeleteOptions{}); err != nil {
		return utils.Fail("系统内部错误"), err
	}
	return utils.Ok("操作成功"), nil
}

// ResetPassword implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !userSrv.CheckUserAllowed(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail("不允许操作超级管理员用户"), nil
	}
	if !userSrv.CheckUserDataScope(ctx, sysUser.UserId, &api.GetOptions{Cache: false}) {
		return utils.Fail("没有权限访问用户数据！"), nil
	}

	sysUser.UpdateBy = loginInfo.User.UserName
	encryptedPasswd, _ := auth.Encrypt(sysUser.Password)
	sysUser.Password = encryptedPasswd
	if err = userSrv.ResetPwd(ctx, sysUser, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// UpdateUserAvatar implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) UpdateUserAvatar(ctx context.Context, req *v1.UpdateUserAvatarRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if err = service.NewUsers(store).UpdateUserAvatar(ctx, loginInfo.User.UserName, req.Avatar, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// ChangeSysUserStatus implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ChangeSysUserStatus(ctx context.Context, req *v1.ChangeSysUserStatus) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	sysUser := v1.UserInfo2SysUser(req.GetUserInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !userSrv.CheckUserAllowed(ctx, sysUser, &api.GetOptions{Cache: false}) {
		return utils.Fail("不允许操作超级管理员用户"), nil
	}
	if !userSrv.CheckUserDataScope(ctx, sysUser.UserId, &api.GetOptions{Cache: false}) {
		return utils.Fail("没有权限访问用户数据！"), nil
	}

	sysUser.UpdateBy = loginInfo.User.UserName
	if err = userSrv.UpdateUserStatus(ctx, sysUser, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), nil
	}
	return utils.Ok("操作成功"), nil
}

// GetAuthRoleById implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) GetAuthRoleById(ctx context.Context, req *v1.GetAuthRoleByIdRequest) (resp *v1.AuthRoleInfoResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	user := service.NewUsers(store).SelectUserById(ctx, req.Id, &api.GetOptions{Cache: true})
	roles := service.NewRoles(store).SelectRolesByUserId(ctx, req.Id, &api.GetOptions{Cache: true})

	return &v1.AuthRoleInfoResponse{
		BaseResp: utils.Ok("操作成功"),
		User:     v1.SysUser2UserInfo(user),
		Roles:    v1.MSysRole2RoleInfo(roles.Items),
	}, nil
}

// AuthRole implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) AuthRole(ctx context.Context, req *v1.AuthRoleRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	userSrv := service.NewUsers(store)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)

	if !userSrv.CheckUserDataScope(ctx, req.UserId, &api.GetOptions{Cache: false}) {
		return utils.Fail("没有权限访问用户数据！"), nil
	}
	if err = userSrv.InsertUserAuth(ctx, req.UserId, req.RoleIds, &api.UpdateOptions{}); err != nil {
		return utils.Fail("系统内部错误"), err
	}

	return utils.Ok("操作成功"), nil
}

// ListDeptsTree implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListDeptsTree(ctx context.Context, req *v1.ListDeptsTreeRequest) (resp *v1.ListDeptsTreeResponse, err error) {
	store, _ := es.GetESFactoryOr(nil)
	dept := v1.DeptInfo2SysDept(req.GetDeptInfo())
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	list := service.NewDepts(store).SelectDeptTreeList(ctx, dept, &api.GetOptions{Cache: true})

	return &v1.ListDeptsTreeResponse{
		BaseResp: utils.Ok("操作成功"),
		Data:     v1.MTreeSelectTrans(list),
	}, nil
}

// ListSysUserOnlines implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ListSysUserOnlines(ctx context.Context, req *v1.ListSysUserOnlinesRequest) (resp *v1.ListSysUserOnline, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	log.Infof("%v %v", req.Ipaddr, req.UserName)
	list := service.NewUserOnlines(store).SelectUserOnline(ctx, req.Ipaddr, req.UserName, &api.GetOptions{Cache: true})
	return &v1.ListSysUserOnline{
		BaseResp: utils.Ok("操作成功"),
		Total:    int64(len(list)),
		Rows:     list,
	}, nil
}

// ForceLogout implements the SystemServiceImpl interface.
func (s *SystemServiceImpl) ForceLogout(ctx context.Context, req *v1.ForceLogoutRequest) (resp *v1.BaseResp, err error) {
	store, _ := es.GetESFactoryOr(nil)
	loginInfo := req.User
	ctx = context.WithValue(ctx, api.LOGIN_INFO_KEY, loginInfo)
	service.NewUserOnlines(store).ForceLogout(ctx, req.TokenId)
	return utils.Ok("操作成功"), nil
}
