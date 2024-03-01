package mysql

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/gateway/v1"
	v12 "github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/pkg/db/kv/redis"
	"github.com/user823/Sophie/pkg/db/sql"
	"testing"
	"time"
)

var (
	ctx context.Context
)

func testInit() {
	cfg := &sql.MysqlConfig{
		Host:                  "127.0.0.1:3306",
		Username:              "sophie",
		Password:              "123456",
		Database:              "sophie",
		MaxIdleConnections:    10,
		MaxOpenConnections:    10,
		MaxConnectionLifeTime: 3600 * time.Second,
		LogLevel:              2,
	}
	_, err := GetMySQLFactoryOr(cfg)
	if err != nil {
		fmt.Printf("出错了 %s", err.Error())
		panic(err)
	}

	testLogininfo := v1.LoginUser{
		User: v12.SysUser{
			UserId:      2,
			DeptId:      105,
			Username:    "sophie",
			Nickname:    "sophie",
			Email:       "sophie@qq.com",
			Phonenumber: "15666666666",
			Sex:         "1",
			Password:    "$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2",
			DelFlag:     "0",
			Roles: []v12.SysRole{
				{RoleId: 2, RoleName: "普通角色", RoleKey: "common", RoleSort: 2, DataScope: "2", MenuCheckStrictly: true, DeptCheckStrictly: true},
			},
		},
		Roles:       []string{"common"},
		Permissions: []string{"*.*.*"},
	}
	ctx = context.WithValue(context.Background(), api.LOGIN_INFO_KEY, testLogininfo)

	connectionConfig := &redis.RedisConfig{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "123456",
		Database: 0,
	}

	go redis.KeepConnection(ctx, connectionConfig)
	time.Sleep(2 * time.Second)
	if !redis.Connected() {
		fmt.Printf("redis 未连接成功")
	}
}

/*
Test SysUser
*/
func TestSelectUserList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	getOpt := &api.GetOptions{
		PageNum:       1,
		PageSize:      10,
		OrderByColumn: "dept_id",
		IsAsc:         false,

		EndTime: time.Now().Unix(),
	}
	sysUser := &v12.SysUser{}
	result, err := sqlCli.Users().SelectUserList(ctx, sysUser, getOpt)
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}

	t.Log("------")
	// 打印角色信息
	for i := range result {
		t.Logf("%v", result[i].Roles)
	}
}

func TestSelectAllocatedList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Users().SelectAllocatedList(ctx, &v12.SysUser{RoleId: 2}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectUnallocatedList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Users().SelectUnallocatedList(ctx, &v12.SysUser{RoleId: 2}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectUserByUserName(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Users().SelectUserByUserName(ctx, "sophie", &api.GetOptions{Cache: true})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestSelectUserById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Users().SelectUserById(ctx, 1, &api.GetOptions{Cache: true})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestCheckUserNameUnique(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result := sqlCli.Users().CheckUserNameUnique(ctx, "sophie", &api.GetOptions{Cache: true})
	t.Logf("%v", result)
}

func TestInsertUser(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	user := &v12.SysUser{
		Username: "test",
		Nickname: "test",
	}
	err := sqlCli.Users().InsertUser(ctx, user, &api.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Users().UpdateUser(ctx, &v12.SysUser{
		UserId: 1,
		Avatar: "123",
	}, &api.UpdateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUserAvatar(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Users().UpdateUserAvatar(ctx, "sophie", "www.baidu.com", &api.UpdateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteUserById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Users().DeleteUserById(ctx, 102, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteUserByIds(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Users().DeleteUserByIds(ctx, []int64{103}, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
Test SysPost
*/
func TestSelectPostList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Posts().SelectPostList(ctx, &v12.SysPost{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectPostById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Posts().SelectPostById(ctx, 1, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestSelectPostListByUserId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Posts().SelectPostListByUserId(ctx, 1, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectPostsByUserName(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Posts().SelectPostsByUserName(ctx, "sophie", &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestUpdatePost(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Posts().UpdatePost(ctx, &v12.SysPost{PostId: 1, PostName: "董事长0"}, &api.UpdateOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
Test SysRole
*/
func TestSelectRoleList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Roles().SelectRoleList(ctx, &v12.SysRole{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectRoleListByUserId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Roles().SelectRoleListByUserId(ctx, 1, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectRoleById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Roles().SelectRoleById(ctx, 2, &api.GetOptions{Cache: false})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestSelectRolesByUserName(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Roles().SelectRolesByUserName(ctx, "sophie", &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestCheckRoleNameUnique(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result := sqlCli.Roles().CheckRoleNameUnique(ctx, "普通角色", &api.GetOptions{Cache: false})
	if result != nil {
		t.Logf("%v", result)
	}
}

func TestInsertRole(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Roles().InsertRole(ctx, &v12.SysRole{RoleName: "test"}, &api.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteRoleById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Roles().DeleteRoleById(ctx, 100, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
	Test SysConfig
*/

func TestSelectConfig(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Configs().SelectConfig(ctx, &v12.SysConfig{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestSelectConfigById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Configs().SelectConfigById(ctx, 1, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestSelectConfigList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Configs().SelectConfigList(ctx, &v12.SysConfig{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

/*
Test SysDept
*/
func TestSelectDeptList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Depts().SelectDeptList(ctx, &v12.SysDept{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectDeptListByRoleId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Depts().SelectDeptListByRoleId(ctx, 2, true, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectChildrenDeptById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Depts().SelectChildrenDeptById(ctx, 100, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectNormalChildrenDeptById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result := sqlCli.Depts().SelectNormalChildrenDeptById(ctx, 100, &api.GetOptions{})
	t.Logf("%v", result)
}

func TestUpdateDeptStatusNormal(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Depts().UpdateDeptStatusNormal(ctx, []int64{100, 105}, &api.UpdateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDeptChildren(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	depts := []*v12.SysDept{
		{DeptId: 200, Ancestors: "0"},
		{DeptId: 201, Ancestors: "1"},
		{DeptId: 202, Ancestors: "2"},
	}
	err := sqlCli.Depts().UpdateDeptChildren(ctx, depts, &api.UpdateOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
	Test SysDictData
*/

func TestSelectDictDataList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.DictData().SelectDictDataList(ctx, &v12.SysDictData{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectDictLabel(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.DictData().SelectDictLabel(ctx, "sys_user_sex", "0", &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestCountDictDataByType(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result := sqlCli.DictData().CountDictDataByType(ctx, "sys_user_sex", &api.GetOptions{})
	t.Log(result)
}

func TestInsertDictData(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	dictData := &v12.SysDictData{
		DictCode:  40,
		DictLabel: "test",
	}
	err := sqlCli.DictData().InsertDictData(ctx, dictData, &api.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteDictDataByIds(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.DictData().DeleteDictDataByIds(ctx, []int64{40}, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
Test DictType
*/
func TestSelectDictTypeList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.DictTypes().SelectDictTypeList(ctx, &v12.SysDictType{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectDictTypeByType(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.DictTypes().SelectDictTypeByType(ctx, "sys_user_sex", &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

/*
	Test SysLogininfor
*/

func TestCleanLogininfor(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Logininfors().CleanLogininfor(ctx, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestSelectLogininforList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Logininfors().SelectLogininforList(ctx, &v12.SysLogininfor{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

/*
	Test SysMenu
*/

func TestSelectMenuList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuList(ctx, &v12.SysMenu{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuPerms(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuPerms(ctx, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuListByUserId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuListByUserId(ctx, &v12.SysMenu{}, 2, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuPermsByRoleId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuPermsByRoleId(ctx, 2, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuPermsByUserId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuPermsByUserId(ctx, 2, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuTreeAll(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuTreeAll(ctx, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuTreeByUserId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuTreeByUserId(ctx, 2, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectMenuListByRoleId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Menus().SelectMenuListByRoleId(ctx, 2, true, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

/*
	test SysNotice
*/

func TestSelectNoticeById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Notices().SelectNoticeById(ctx, 1, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestSelectNoticeList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.Notices().SelectNoticeList(ctx, &v12.SysNotice{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestInsertNotice(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.Notices().InsertNotice(ctx, &v12.SysNotice{NoticeType: "C", NoticeTitle: "test", NoticeContent: "test"}, &api.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
	Test SysOperLog
*/

func TestInsertOperLog(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.OperLogs().InsertOperLog(ctx, &v12.SysOperLog{Title: "test", OperTime: time.Now()}, &api.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestSelectOperLogList(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.OperLogs().SelectOperLogList(ctx, &v12.SysOperLog{}, &api.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	for i := range result {
		t.Logf("%v", result[i])
	}
}

func TestSelectOperLogById(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	result, err := sqlCli.OperLogs().SelectOperLogById(ctx, 101, &api.GetOptions{Cache: true})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", result)
}

func TestCleanOperLog(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.OperLogs().CleanOperLog(ctx, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

/*
	Test SysRoleDept
*/

func TestBatchRoleDept(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	list := []*v12.SysRoleDept{
		{RoleId: 3, DeptId: 4},
		{RoleId: 5, DeptId: 6},
	}
	err := sqlCli.RoleDepts().BatchRoleDept(ctx, list, &api.CreateOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestSelectCountRoleDeptByDeptId(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	c := sqlCli.RoleDepts().SelectCountRoleDeptByDeptId(ctx, 101, &api.GetOptions{})
	t.Logf("%v", c)
}

func TestDeleteRoleDept(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	err := sqlCli.RoleDepts().DeleteRoleDept(ctx, []int64{3, 5}, &api.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}
}

func TestTxRollback(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	tx := sqlCli.Begin()
	tx.Users().InsertUser(ctx, &v12.SysUser{UserId: 5, Username: "test"}, &api.CreateOptions{})
	tx.Roles().InsertRole(ctx, &v12.SysRole{RoleId: 5, RoleKey: "test", RoleName: "test"}, &api.CreateOptions{})
	tx.Rollback()
}

func TestTxCommit(t *testing.T) {
	sqlCli, _ := GetMySQLFactoryOr(nil)
	tx := sqlCli.Begin()
	tx.Users().InsertUser(ctx, &v12.SysUser{UserId: 5, Username: "test"}, &api.CreateOptions{})
	tx.Roles().InsertRole(ctx, &v12.SysRole{RoleId: 5, RoleKey: "test", RoleName: "test"}, &api.CreateOptions{})
	tx.Commit()
}

func TestSub(t *testing.T) {
	testInit()

	t.Run("test-SelectUserList", TestSelectUserList)
	t.Run("test-SelectAllocatedList", TestSelectAllocatedList)
	t.Run("test-SelectUnallocatedList", TestSelectUnallocatedList)
	t.Run("test-SelectUserByUserName", TestSelectUserByUserName)
	t.Run("test-SelectUserById", TestSelectUserById)
	t.Run("test-CheckUserNameUnique", TestCheckUserNameUnique)
	t.Run("test-InsertUser", TestInsertUser)
	t.Run("test-UpdateUser", TestUpdateUser)
	t.Run("test-UpdateUserAvatar", TestUpdateUserAvatar)
	t.Run("test-DeleteUserById", TestDeleteUserById)
	t.Run("test-DeleteUserByIds", TestDeleteUserByIds)

	t.Run("test-SelectPostList", TestSelectPostList)
	t.Run("test-SelectPostById", TestSelectPostById)
	t.Run("test-SelectPostListByUserId", TestSelectPostListByUserId)
	t.Run("test-SelectPostsByUserName", TestSelectPostsByUserName)
	t.Run("test-UpdatePost", TestUpdatePost)

	t.Run("test-SelectRoleList", TestSelectRoleList)
	t.Run("test-SelectRoleListByUserId", TestSelectRoleListByUserId)
	t.Run("test-SelectRoleById", TestSelectRoleById)
	t.Run("test-SelectRolesByUserName", TestSelectRolesByUserName)
	t.Run("test-CheckRoleNameUnique", TestCheckRoleNameUnique)
	t.Run("test-InsertRole", TestInsertRole)
	t.Run("test-DeleteRoleById", TestDeleteRoleById)

	t.Run("test-SelectConfig", TestSelectConfig)
	t.Run("test-SelectConfigById", TestSelectConfigById)
	t.Run("test-SelectConfigList", TestSelectConfigList)

	t.Run("test-SelectDeptList", TestSelectDeptList)
	t.Run("test-SelectDeptListByRoleId", TestSelectDeptListByRoleId)
	t.Run("test-SelectChildrenDeptById", TestSelectChildrenDeptById)
	t.Run("test-SelectNormalChildrenDeptById", TestSelectNormalChildrenDeptById)
	t.Run("test-UpdateDeptStatusNormal", TestUpdateDeptStatusNormal)
	t.Run("test-UpdateDeptChildren", TestUpdateDeptChildren)

	t.Run("test-SelectDictDataList", TestSelectDictDataList)
	t.Run("test-SelectDictLabel", TestSelectDictLabel)
	t.Run("test-CountDictDataByType", TestCountDictDataByType)
	t.Run("test-InsertDictData", TestInsertDictData)
	t.Run("test-DeleteDictDataByIds", TestDeleteDictDataByIds)

	t.Run("test-SelectDictTypeList", TestSelectDictTypeList)
	t.Run("test-SelectDictTypeByType", TestSelectDictTypeByType)

	t.Run("test-SelectLogininforList", TestSelectLogininforList)
	t.Run("test-CleanLogininfor", TestCleanLogininfor)

	t.Run("test-SelectMenuList", TestSelectMenuList)
	t.Run("test-SelectMenuPerms", TestSelectMenuPerms)
	t.Run("test-SelectMenuListByUserId", TestSelectMenuListByUserId)
	t.Run("test-SelectMenuPermsByRoleId", TestSelectMenuPermsByRoleId)
	t.Run("test-SelectMenuPermsByUserId", TestSelectMenuPermsByUserId)
	t.Run("test-SelectMenuTreeAll", TestSelectMenuTreeAll)
	t.Run("test-SelectMenuTreeByUserId", TestSelectMenuTreeByUserId)
	t.Run("test-SelectMenuListByRoleId", TestSelectMenuListByRoleId)

	t.Run("test-SelectNoticeById", TestSelectNoticeById)
	t.Run("test-SelectNoticeList", TestSelectNoticeList)
	t.Run("test-InsertNotice", TestInsertNotice)

	t.Run("test-InsertOperLog", TestInsertOperLog)
	t.Run("test-SelectOperLogList", TestSelectOperLogList)
	t.Run("test-SelectOperLogById", TestSelectOperLogById)
	t.Run("test-CleanOperLog", TestCleanOperLog)

	t.Run("test-BatchRoleDept", TestBatchRoleDept)
	t.Run("test-SelectCountRoleDeptByDeptId", TestSelectCountRoleDeptByDeptId)
	t.Run("test-DeleteRoleDept", TestDeleteRoleDept)

	t.Run("test-TxRollback", TestTxRollback)
	t.Run("test-TxCommit", TestTxCommit)
}

// SELECT `u`.`create_by`,`u`.`create_time`,`u`.`update_by`,`u`.`update_time`,`u`.`remark`,`u`.`extend_shadow`,`u`.`user_id`,`u`.`dept_id`,`u`.`user_name`,`u`.`nick_name`,`u`.`email`,`u`.`phonenumber`,`u`.`sex`,`u`.`avatar`,`u`.`password`,`u`.`status`,`u`.`del_flag`,`u`.`login_ip`,`u`.`login_date` FROM sys_user u left join sys_dept d on u.dept_id = d.dept_id left join sys_user_role ur on u.user_id = ur.user_id left join sys_role r on r.role_id = ur.role_id WHERE u.del_flag = 0 AND u.create_time <= '2024-02-21 08:26:41' AND d.dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id = 2 ORDER BY dept_id DESC LIMIT 10
