package service

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/user823/Sophie/api"
	"github.com/user823/Sophie/api/domain/system/v1"
	vo2 "github.com/user823/Sophie/api/domain/vo"
	"github.com/user823/Sophie/internal/system/store"
	"github.com/user823/Sophie/pkg/utils/intutil"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strings"
)

type MenuSrv interface {
	// 根据用户查询系统菜单列表
	SelectMenuList(ctx context.Context, userId int64, opts *api.GetOptions) *v1.MenuList
	// 根据用户查询系统菜单列表
	SelectMenuListWithMenu(ctx context.Context, menu *v1.SysMenu, userId int64, opts *api.GetOptions) *v1.MenuList
	// 根据用户ID查询权限
	SelectMenuPermsByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []string
	// 根据角色ID查询权限
	SelectMenuPermsByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) []string
	// 根据用户ID查询菜单树信息
	SelectMenuTreeByUserId(ctx context.Context, userId int64, opts *api.GetOptions) *v1.MenuList
	// 根据角色ID查询菜单树信息
	SelectMenuListByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) []int64
	// 构建前端路由所需要的菜单
	BuildMenus(ctx context.Context, menus []*v1.SysMenu) []vo2.RouterVo
	// 构建前端路由所需要的树结构
	BuildMenuTree(ctx context.Context, menus []*v1.SysMenu) *v1.MenuList
	// 构建前端所需要下拉树结构
	BuildMenuTreeSelect(ctx context.Context, menus []*v1.SysMenu) []vo2.TreeSelect
	// 根据菜单ID查询信息
	SelectMenuById(ctx context.Context, menuId int64, opts *api.GetOptions) *v1.SysMenu
	// 是否存在菜单子节点
	HasChildByMenuId(ctx context.Context, menuId int64, opts *api.GetOptions) bool
	// 查询菜单是否存在角色
	CheckMenuExistsRole(ctx context.Context, menuId int64, opts *api.GetOptions) bool
	// 新增保存菜单信息
	InsertMenu(ctx context.Context, menu *v1.SysMenu, opts *api.CreateOptions) error
	// 修改保存菜单信息
	UpdateMenu(ctx context.Context, menu *v1.SysMenu, opts *api.UpdateOptions) error
	// 删除菜单管理信息
	DeleteMenuBuId(ctx context.Context, menuId int64, opts *api.DeleteOptions) error
	// 校验菜单名称是否唯一
	CheckMenuNameUnique(ctx context.Context, menu *v1.SysMenu, opts *api.GetOptions) bool
}

type menuService struct {
	store store.Factory
}

var _ MenuSrv = &menuService{}

func NewMenus(s store.Factory) MenuSrv {
	return &menuService{s}
}

func (s *menuService) SelectMenuList(ctx context.Context, userId int64, opts *api.GetOptions) *v1.MenuList {
	return s.SelectMenuListWithMenu(ctx, &v1.SysMenu{}, userId, opts)
}

func (s *menuService) SelectMenuListWithMenu(ctx context.Context, menu *v1.SysMenu, userId int64, opts *api.GetOptions) *v1.MenuList {
	var result []*v1.SysMenu
	var err error

	// 管理员搜索所有列表
	if v1.IsUserAdmin(userId) {
		result, err = s.store.Menus().SelectMenuList(ctx, menu, opts)
		if err != nil {
			return &v1.MenuList{ListMeta: api.ListMeta{0}}
		}
	} else {
		result, err = s.store.Menus().SelectMenuListByUserId(ctx, menu, userId, opts)
		if err != nil {
			return &v1.MenuList{ListMeta: api.ListMeta{0}}
		}
	}

	return &v1.MenuList{ListMeta: api.ListMeta{int64(len(result))}, Items: result}
}

func (s *menuService) SelectMenuPermsByUserId(ctx context.Context, userId int64, opts *api.GetOptions) []string {
	perms, err := s.store.Menus().SelectMenuPermsByUserId(ctx, userId, opts)
	if err != nil {
		return []string{}
	}
	permSet := hashset.New()
	for _, perm := range perms {
		if perm != "" {
			for _, str := range strings.Split(strings.TrimSpace(perm), ",") {
				permSet.Add(str)
			}
		}
	}

	res := make([]string, 0, permSet.Size())
	for _, v := range permSet.Values() {
		res = append(res, v.(string))
	}
	return res
}

func (s *menuService) SelectMenuPermsByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) []string {
	perms, err := s.store.Menus().SelectMenuPermsByRoleId(ctx, roleId, opts)
	if err != nil {
		return []string{}
	}
	permSet := hashset.New()
	for _, perm := range perms {
		if perm != "" {
			for _, str := range strings.Split(strings.TrimSpace(perm), ",") {
				permSet.Add(str)
			}
		}
	}

	res := make([]string, 0, permSet.Size())
	for _, v := range permSet.Values() {
		res = append(res, v.(string))
	}
	return res
}

func (s *menuService) SelectMenuTreeByUserId(ctx context.Context, userId int64, opts *api.GetOptions) *v1.MenuList {
	var result []*v1.SysMenu
	var err error

	// 如果是超级用户则获部门信息
	if v1.IsUserAdmin(userId) {
		result, err = s.store.Menus().SelectMenuTreeAll(ctx, opts)
		if err != nil {
			return &v1.MenuList{ListMeta: api.ListMeta{0}}
		}
	} else {
		result, err = s.store.Menus().SelectMenuTreeByUserId(ctx, userId, opts)
		if err != nil {
			return &v1.MenuList{ListMeta: api.ListMeta{0}}
		}
	}

	root := &v1.SysMenu{MenuId: 0}
	menuRecursionFn(result, root)
	return &v1.MenuList{ListMeta: api.ListMeta{int64(len(root.Children))}, Items: root.Children}
}

func (s *menuService) SelectMenuListByRoleId(ctx context.Context, roleId int64, opts *api.GetOptions) []int64 {
	role, err := s.store.Roles().SelectRoleById(ctx, roleId, opts)
	if err != nil {
		return []int64{}
	}
	result, err := s.store.Menus().SelectMenuListByRoleId(ctx, roleId, role.MenuCheckStrictly, opts)
	if err != nil {
		return []int64{}
	}
	return result
}

func (s *menuService) BuildMenus(ctx context.Context, menus []*v1.SysMenu) []vo2.RouterVo {
	routers := make([]vo2.RouterVo, 0, len(menus))
	for i := range menus {
		router := vo2.RouterVo{}
		router.Hidden = menus[i].Visible == "1"
		router.Name = getRouterName(menus[i])
		router.Path = getRouterPath(menus[i])
		router.Component = getComponent(menus[i])
		router.Query = menus[i].Query
		router.Meta = NewMetaVo(menus[i].MenuName, menus[i].Icon, menus[i].IsCache == "1", menus[i].Path)

		cMenus := menus[i].Children
		if len(cMenus) > 0 && v1.TYPE_DIR == menus[i].MenuType {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
			router.Children = s.BuildMenus(ctx, cMenus)
		} else if isMenuFrame(menus[i]) {
			router.Meta = vo2.MetaVo{}
			children := vo2.RouterVo{}
			children.Path = menus[i].Path
			children.Component = menus[i].Component
			children.Name = strutil.Capitalize(menus[i].Path)
			children.Meta = NewMetaVo(menus[i].MenuName, menus[i].Icon, menus[i].IsCache == "1", menus[i].Path)
			children.Query = menus[i].Query
			router.Children = []vo2.RouterVo{children}
		} else if menus[i].ParentId == 0 && isInnerLink(menus[i]) {
			router.Meta = vo2.MetaVo{menus[i].MenuName, menus[i].Icon, false, ""}
			router.Path = "/"
			children := vo2.RouterVo{}
			routerPath := innerLinkReplaceEach(menus[i].Path)
			children.Path = routerPath
			children.Component = v1.INNER_LINK
			children.Name = strutil.Capitalize(routerPath)
			children.Meta = vo2.MetaVo{menus[i].MenuName, menus[i].Icon, false, menus[i].Path}
			router.Children = []vo2.RouterVo{children}
		}
		routers = append(routers, router)
	}
	return routers
}

func (s *menuService) BuildMenuTree(ctx context.Context, menus []*v1.SysMenu) *v1.MenuList {
	list := make([]*v1.SysMenu, 0, len(menus))
	menuIds := make([]int64, 0, len(menus))
	for i := range menus {
		menuIds = append(menuIds, menus[i].MenuId)
	}

	for i := range menus {
		// 顶层节点，则递归
		if !intutil.ContainsAnyInt64(menus[i].ParentId, menuIds...) {
			menuRecursionFn(menus, menus[i])
			list = append(list, menus[i])
		}
	}
	if len(list) == 0 {
		list = menus
	}
	return &v1.MenuList{
		ListMeta: api.ListMeta{int64(len(list))},
		Items:    list,
	}
}

func (s *menuService) BuildMenuTreeSelect(ctx context.Context, menus []*v1.SysMenu) []vo2.TreeSelect {
	menuTrees := s.BuildMenuTree(ctx, menus)
	res := make([]vo2.TreeSelect, 0, menuTrees.ListMeta.TotalCount)
	for i := range menuTrees.Items {
		res = append(res, menuTrees.Items[i].BuildTreeSelect())
	}
	return res
}

func (s *menuService) SelectMenuById(ctx context.Context, menuId int64, opts *api.GetOptions) *v1.SysMenu {
	result, err := s.store.Menus().SelectMenuById(ctx, menuId, opts)
	if err != nil {
		return nil
	}
	return result
}

func (s *menuService) HasChildByMenuId(ctx context.Context, menuId int64, opts *api.GetOptions) bool {
	return s.store.Menus().HasChildByMenuId(ctx, menuId, opts)
}

func (s *menuService) CheckMenuExistsRole(ctx context.Context, menuId int64, opts *api.GetOptions) bool {
	return s.store.RoleMenus().CheckMenuExistRole(ctx, menuId, opts) > 0
}

func (s *menuService) InsertMenu(ctx context.Context, menu *v1.SysMenu, opts *api.CreateOptions) error {
	return s.store.Menus().InsertMenu(ctx, menu, opts)
}

func (s *menuService) UpdateMenu(ctx context.Context, menu *v1.SysMenu, opts *api.UpdateOptions) error {
	return s.store.Menus().UpdateMenu(ctx, menu, opts)
}

func (s *menuService) DeleteMenuBuId(ctx context.Context, menuId int64, opts *api.DeleteOptions) error {
	return s.store.Menus().DeleteMenuById(ctx, menuId, opts)
}

func (s *menuService) CheckMenuNameUnique(ctx context.Context, menu *v1.SysMenu, opts *api.GetOptions) bool {
	info := s.store.Menus().CheckMenuNameUnique(ctx, menu.MenuName, menu.ParentId, opts)
	if info != nil && info.MenuId != menu.MenuId {
		return false
	}
	return true
}

func getRouterName(menu *v1.SysMenu) string {
	routerName := strutil.Capitalize(menu.Path)
	// 非外链并且是一级目录（类型为目录）
	if isMenuFrame(menu) {
		routerName = ""
	}
	return routerName
}

func getRouterPath(menu *v1.SysMenu) string {
	routerPath := menu.Path
	// 内链打开外网方式
	if menu.ParentId != 0 && isInnerLink(menu) {
		routerPath = innerLinkReplaceEach(routerPath)
	}
	// 非外链并且是一级目录（类型为目录）
	if menu.ParentId == 0 && v1.TYPE_DIR == menu.MenuType && v1.NO_FRAME == menu.IsFrame {
		routerPath = "/" + menu.Path
	} else if isMenuFrame(menu) { // 非外链并且是一级目录（类型为菜单）
		routerPath = "/"
	}
	return routerPath
}

// 非外链并且是一级目录（类型为菜单）
func isMenuFrame(menu *v1.SysMenu) bool {
	return menu.ParentId == 0 && v1.TYPE_MENU == menu.MenuType && v1.NO_FRAME == menu.IsFrame
}

// 是否为内链
func isInnerLink(menu *v1.SysMenu) bool {
	return menu.IsFrame == v1.NO_FRAME && strutil.IsHttp(menu.Path)
}

// 内链域名特殊字符替换
func innerLinkReplaceEach(path string) string {
	oldStr := []string{api.HTTP, api.HTTPS, api.WWW, ".", ":"}
	newStr := []string{"", "", "", "/", "/"}
	return strutil.ReplaceEach(path, oldStr, newStr)
}

// 获取组件信息
func getComponent(menu *v1.SysMenu) string {
	component := v1.LAYOUT
	if menu.Component != "" && !isMenuFrame(menu) {
		component = menu.Component
	} else if menu.Component != "" && menu.ParentId != 0 && isInnerLink(menu) {
		component = v1.INNER_LINK
	} else if menu.Component == "" && isParentView(menu) {
		component = v1.PARENT_VIEW
	}
	return component
}

// 是否为parent_view组件
func isParentView(menu *v1.SysMenu) bool {
	return menu.ParentId != 0 && v1.TYPE_DIR == menu.MenuType
}

// 递归列表
func menuRecursionFn(list []*v1.SysMenu, t *v1.SysMenu) {
	// 子节点列表
	childList := getMenuChildList(list, t)
	t.Children = childList
	for i := range childList {
		if len(getMenuChildList(list, childList[i])) > 0 {
			menuRecursionFn(list, childList[i])
		}
	}
}

// 得到子节点列表
func getMenuChildList(list []*v1.SysMenu, t *v1.SysMenu) []*v1.SysMenu {
	tlist := make([]*v1.SysMenu, 0, len(list))
	for i := range list {
		if list[i].ParentId == t.MenuId {
			tlist = append(tlist, list[i])
		}
	}
	return tlist
}

func NewMetaVo(title, icon string, noCache bool, link string) vo2.MetaVo {
	res := vo2.MetaVo{
		Title:   title,
		Icon:    icon,
		NoCache: noCache,
	}
	if strutil.IsHttp(link) {
		res.Link = link
	}
	return res
}
