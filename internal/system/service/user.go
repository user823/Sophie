package service

import (
	"context"
	"fmt"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	"github.com/user823/Sophie/internal/pkg/middleware/auth"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/errors"
	"strings"
)

type UserSrv interface {
	// 根据分页条件查询用户列表
	SelectUserList(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.UserList
	// 根据条件分页查询已分配用户角色列表
	SelectAllocatedList(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.UserList
	// 根据分页条件查询未分配用户角色列表
	SelectUnallocatedList(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.UserList
	// 通过用户名查询用户
	SelectUserByUserName(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.SysUser
	// 通过用户ID查询用户
	SelectUserById(ctx context.Context, id int64, opts *api.GetOptions) *v1.SysUser
	// 根据用户ID查询用户所属角色组
	SelectUserRoleGroup(ctx context.Context, userName string, opts *api.GetOptions) string
	// 根据用户id查询用户所属岗位组
	SelectUserPostGroup(ctx context.Context, userName string, opts *api.GetOptions) string
	// 校验用户名称是否唯一
	CheckUserNameUnique(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool
	// 校验手机号是否唯一
	CheckPhoneUnique(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool
	// 检验email是否唯一
	CheckEmailUnique(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool
	// 校验用户是否允许操作
	CheckUserAllowed(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool
	// 校验用户是否有数据权限
	CheckUserDataScope(ctx context.Context, id int64, opts *api.GetOptions) bool
	// 新增用户信息
	InsertUser(ctx context.Context, user *v1.SysUser, opts *api.CreateOptions) error
	// 注册用户信息
	RegisterUser(ctx context.Context, user *v1.SysUser, opts *api.CreateOptions) bool
	// 修改用户信息
	UpdateUser(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error
	// 用户授权角色
	InsertUserAuth(ctx context.Context, id int64, roleIds []int64, opts *api.UpdateOptions) error
	// 修改用户状态
	UpdateUserStatus(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error
	// 修改用户基本信息
	UpdateUsrProfile(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error
	// 修改用户头像
	UpdateUserAvatar(ctx context.Context, userName, avatar string, opts *api.UpdateOptions) error
	// 重置用户密码
	ResetPwd(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error
	// 重置用户密码
	ResetUserPwd(ctx context.Context, userName, password string, opts *api.UpdateOptions) error
	// 通过用户ID删除用户
	DeleteUserById(ctx context.Context, userid int64, opts *api.DeleteOptions) error
	// 批量删除用户信息
	DeleteUserByIds(ctx context.Context, userIds []int64, opts *api.DeleteOptions) error
	// 导入用户数据
	ImportUser(ctx context.Context, userList []*v1.SysUser, isUpdateSupport bool, operName string, opts *api.UpdateOptions) (string, error)
}

type userService struct {
	store store.Factory
}

var _ UserSrv = &userService{}

var (
	ErrUserNotAllowed = fmt.Errorf("不允许操作超级管理员用户")
	ErrUserDataScope  = fmt.Errorf("没有权限访问用户数据")
)

func NewUsers(s store.Factory) UserSrv {
	return &userService{store: s}
}

func (s *userService) SelectUserList(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.UserList {
	result, _ := s.store.Users().SelectUserList(ctx, user, opts)
	return &v1.UserList{
		ListMeta: api.ListMeta{int64(len(result))},
		Items:    result,
	}
}

func (s *userService) SelectAllocatedList(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.UserList {
	result, _ := s.store.Users().SelectAllocatedList(ctx, user, opts)
	return &v1.UserList{
		ListMeta: api.ListMeta{int64(len(result))},
		Items:    result,
	}
}

func (s *userService) SelectUnallocatedList(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.UserList {
	result, _ := s.store.Users().SelectUnallocatedList(ctx, user, opts)
	return &v1.UserList{
		ListMeta: api.ListMeta{int64(len(result))},
		Items:    result,
	}
}

func (s *userService) SelectUserByUserName(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) *v1.SysUser {
	result, _ := s.store.Users().SelectUserByUserName(ctx, user.Username, opts)
	return result
}

func (s *userService) SelectUserById(ctx context.Context, id int64, opts *api.GetOptions) *v1.SysUser {
	result, _ := s.store.Users().SelectUserById(ctx, id, opts)
	return result
}

func (s *userService) SelectUserRoleGroup(ctx context.Context, userName string, opts *api.GetOptions) string {
	roles, _ := s.store.Roles().SelectRolesByUserName(ctx, userName, opts)
	if len(roles) == 0 {
		return ""
	}
	var builder strings.Builder
	for i := range roles {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(roles[i].RoleName)
	}
	return builder.String()
}

func (s *userService) SelectUserPostGroup(ctx context.Context, userName string, opts *api.GetOptions) string {
	posts, err := s.store.Posts().SelectPostsByUserName(ctx, userName, opts)
	if err != nil || len(posts) == 0 {
		return ""
	}

	var builder strings.Builder
	for i := range posts {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(posts[i].PostName)
	}
	return builder.String()
}

func (s *userService) CheckUserNameUnique(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool {
	result := s.store.Users().CheckUserNameUnique(ctx, user.Username, opts)
	if result != nil && user.UserId != result.UserId {
		return false
	}
	return true
}

func (s *userService) CheckPhoneUnique(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool {
	result := s.store.Users().CheckPhoneUnique(ctx, user.Phonenumber, opts)
	if result != nil && user.UserId != result.UserId {
		return false
	}
	return true
}

func (s *userService) CheckEmailUnique(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool {
	result := s.store.Users().CheckEmailUnique(ctx, user.Email, opts)
	if result != nil && user.UserId != result.UserId {
		return false
	}
	return true
}

func (s *userService) CheckUserAllowed(ctx context.Context, user *v1.SysUser, opts *api.GetOptions) bool {
	// 不允许操作超级管理员用户
	return !(user.UserId != 0 && user.IsAdmin())
}

// 判断当前用户能否查到目标用户的信息
func (s *userService) CheckUserDataScope(ctx context.Context, id int64, opts *api.GetOptions) bool {
	logininfor := utils.GetLogininfoFromCtx(ctx)
	if !v1.IsUserAdmin(logininfor.User.UserId) {
		users, err := s.store.Users().SelectUserList(ctx, &v1.SysUser{UserId: id}, opts)
		// 没有权限访问用户数据
		if err != nil || len(users) == 0 {
			return false
		}
	}
	return true
}

func (s *userService) InsertUser(ctx context.Context, user *v1.SysUser, opts *api.CreateOptions) error {
	if err := s.store.Users().InsertUser(ctx, user, opts); err != nil {
		return err
	}
	s.InsertUserPost(ctx, user, &api.CreateOptions{})
	s.InsertUserRole(ctx, user, &api.CreateOptions{})
	return nil
}

func (s *userService) RegisterUser(ctx context.Context, user *v1.SysUser, opts *api.CreateOptions) bool {
	return s.store.Users().InsertUser(ctx, user, opts) == nil
}

func (s *userService) UpdateUser(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error {
	// 删除用户与角色关联
	tx := s.store.Begin()
	if err := tx.UserRoles().DeleteUserRoleByUserId(ctx, user.UserId, &api.DeleteOptions{}); err != nil {
		tx.Rollback()
		return err
	}

	// 新增用户和角色的关联
	if len(user.RoleIds) > 0 {
		list := make([]*v1.SysUserRole, 0, len(user.RoleIds))
		for i := range user.RoleIds {
			list = append(list, &v1.SysUserRole{UserId: user.UserId, RoleId: user.RoleIds[i]})
		}
		if err := tx.UserRoles().BatchUserRole(ctx, list, &api.CreateOptions{}); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除用户与岗位的关联
	if err := tx.UserPosts().DeleteUserPostByUserId(ctx, user.UserId, &api.DeleteOptions{}); err != nil {
		tx.Rollback()
		return err
	}

	// 新增用户与岗位关联
	if len(user.PostIds) > 0 {
		list := make([]*v1.SysUserPost, 0, len(user.PostIds))
		for i := range user.PostIds {
			list = append(list, &v1.SysUserPost{UserId: user.UserId, PostId: user.PostIds[i]})
		}
		if err := tx.UserPosts().BatchUserPost(ctx, list, &api.CreateOptions{}); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新用户信息
	if err := tx.Users().UpdateUser(ctx, user, &api.UpdateOptions{}); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *userService) InsertUserAuth(ctx context.Context, id int64, roleIds []int64, opts *api.UpdateOptions) error {
	tx := s.store.Begin()
	if err := tx.UserRoles().DeleteUserRoleByUserId(ctx, id, &api.DeleteOptions{}); err != nil {
		tx.Rollback()
		return err
	}

	if len(roleIds) > 0 {
		list := make([]*v1.SysUserRole, 0, len(roleIds))
		for i := range roleIds {
			list = append(list, &v1.SysUserRole{UserId: id, RoleId: roleIds[i]})
		}
		if err := tx.UserRoles().BatchUserRole(ctx, list, &api.CreateOptions{}); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *userService) UpdateUserStatus(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error {
	return s.store.Users().UpdateUser(ctx, user, opts)
}

func (s *userService) UpdateUsrProfile(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error {
	return s.store.Users().UpdateUser(ctx, user, opts)
}

func (s *userService) UpdateUserAvatar(ctx context.Context, userName, avatar string, opts *api.UpdateOptions) error {
	return s.store.Users().UpdateUserAvatar(ctx, userName, avatar, opts)
}

func (s *userService) ResetPwd(ctx context.Context, user *v1.SysUser, opts *api.UpdateOptions) error {
	return s.store.Users().UpdateUser(ctx, user, opts)
}

func (s *userService) ResetUserPwd(ctx context.Context, userName, password string, opts *api.UpdateOptions) error {
	return s.store.Users().UpdateUserPwd(ctx, userName, password, opts)
}

func (s *userService) DeleteUserById(ctx context.Context, userid int64, opts *api.DeleteOptions) error {
	tx := s.store.Begin()
	// 删除用户与角色的关联
	if err := tx.UserRoles().DeleteUserRoleByUserId(ctx, userid, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除用户与岗位的关联
	if err := tx.UserPosts().DeleteUserPostByUserId(ctx, userid, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除用户信息
	if err := tx.Users().DeleteUserById(ctx, userid, opts); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *userService) DeleteUserByIds(ctx context.Context, userIds []int64, opts *api.DeleteOptions) error {
	for i := range userIds {
		if !s.CheckUserAllowed(ctx, &v1.SysUser{UserId: userIds[i]}, &api.GetOptions{Cache: true}) {
			return ErrUserNotAllowed
		}
		if !s.CheckUserDataScope(ctx, userIds[i], &api.GetOptions{Cache: true}) {
			return ErrUserDataScope
		}
	}
	tx := s.store.Begin()
	// 删除用户角色关联
	if err := tx.UserRoles().DeleteUserRole(ctx, userIds, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除用户岗位关联
	if err := tx.UserPosts().DeleteUserPost(ctx, userIds, opts); err != nil {
		tx.Rollback()
		return err
	}
	// 删除用户信息
	if err := tx.Users().DeleteUserByIds(ctx, userIds, opts); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *userService) ImportUser(ctx context.Context, userList []*v1.SysUser, isUpdateSupport bool, operName string, opts *api.UpdateOptions) (string, error) {
	if len(userList) == 0 {
		return "", fmt.Errorf("导入用户数据不能为空")
	}
	var successNum, failNum int
	var successMsg, failMsg strings.Builder
	dealFail := func(err error, username string) {
		failNum++
		failMsg.WriteString(fmt.Sprintf("<br/>%d、账号 %s 导入失败：%s", failNum, username, err.Error()))
	}

	password := NewConfigs(s.store).SelectConfigByKey(ctx, "sys.user.initPassword", &api.GetOptions{Cache: true})
	for i := range userList {
		// 验证是否存在用户
		sysUser, _ := s.store.Users().SelectUserByUserName(ctx, userList[i].Username, &api.GetOptions{Cache: true})
		if sysUser == nil {
			// 验证字段合法性
			if err := userList[i].Validate(); err != nil {
				dealFail(err, userList[i].Username)
				continue
			}

			userList[i].Password, _ = auth.Encrypt(password)
			userList[i].CreateBy = operName
			if err := s.store.Users().InsertUser(ctx, userList[i], &api.CreateOptions{}); err != nil {
				dealFail(err, userList[i].Username)
				continue
			}
			successNum++
			successMsg.WriteString(fmt.Sprintf("<br/>%d、账号 %s 导入成功", successNum, userList[i].Username))
		} else if isUpdateSupport {
			if err := userList[i].Validate(); err != nil {
				dealFail(err, userList[i].Username)
				continue
			}
			if !s.CheckUserAllowed(ctx, userList[i], &api.GetOptions{Cache: true}) {
				dealFail(ErrUserNotAllowed, userList[i].Username)
				continue
			}
			if !s.CheckUserDataScope(ctx, userList[i].UserId, &api.GetOptions{Cache: true}) {
				dealFail(ErrUserDataScope, userList[i].Username)
				continue
			}
			userList[i].UpdateBy = operName
			if err := s.store.Users().UpdateUser(ctx, userList[i], opts); err != nil {
				dealFail(err, userList[i].Username)
				continue
			}
			successNum++
			successMsg.WriteString(fmt.Sprintf("<br/>%d、账号 %s 更新成功", successNum, userList[i].Username))
		} else {
			failNum++
			failMsg.WriteString(fmt.Sprintf("<br/>%d、账号 %s 已存在", failNum, userList[i].Username))
		}
	}
	if failNum > 0 {
		msg := fmt.Sprintf("很抱歉，导入失败！共 %d 条数据格式不正确，错误如下：") + failMsg.String()
		return "", errors.New(msg)
	} else {
		msg := fmt.Sprintf("恭喜您，数据已全部导入成功！共 %d 条，数据如下: ") + successMsg.String()
		return msg, nil
	}
}

// 内部使用的方法
func (s *userService) InsertUserPost(ctx context.Context, user *v1.SysUser, opts *api.CreateOptions) {
	if len(user.PostIds) > 0 {
		list := make([]*v1.SysUserPost, 0, len(user.PostIds))
		for i := range user.PostIds {
			list = append(list, &v1.SysUserPost{
				UserId: user.UserId,
				PostId: user.PostIds[i],
			})
		}
		s.store.UserPosts().BatchUserPost(ctx, list, opts)
	}
}

func (s *userService) InsertUserRole(ctx context.Context, user *v1.SysUser, opts *api.CreateOptions) {
	if len(user.RoleIds) > 0 {
		list := make([]*v1.SysUserRole, 0, len(user.RoleIds))
		for i := range user.RoleIds {
			list = append(list, &v1.SysUserRole{
				UserId: user.UserId,
				RoleId: user.RoleIds[i],
			})
		}
		s.store.UserRoles().BatchUserRole(ctx, list, opts)
	}
}
