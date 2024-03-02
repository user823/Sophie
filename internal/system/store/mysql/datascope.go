package mysql

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/user823/Sophie/internal/pkg/middleware/secure"
	"github.com/user823/Sophie/internal/system/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"gorm.io/gorm"
	"strings"
)

const (
	DATA_SCOPE_ALL            = "1"
	DATA_SCOPE_CUSTOM         = "2"
	DATA_SCOPE_DEPT           = "3"
	DATA_SCOPE_DEPT_AND_CHILD = "4"
	DATA_SCOPE_SELF           = "5"
)

// 限制联表查询时的数据范围
func dateScopeFromCtx(ctx context.Context, db *gorm.DB, userAlias string, deptAlias string) (*gorm.DB, error) {
	// 获取登录信息
	logininfo := utils.GetLogininfoFromCtx(ctx)
	// 获取要求的权限数据
	var perm_data string
	data := ctx.Value(secure.REQUIRE_PERMISSION)
	if data != nil {
		perm_data = data.(string)
	}

	conditions := hashset.New()
	roles := logininfo.User.Roles
	user := logininfo.User
	// 要求的权限
	permissions := strings.Split(perm_data, secure.DefaultSeprator)

	var queryCondition strings.Builder

	for _, role := range roles {
		dataScope := role.DataScope
		// 处理过的dataScope，跳过
		if dataScope != DATA_SCOPE_CUSTOM && conditions.Contains(dataScope) {
			continue
		}

		// 判断角色是否满足要求的权限之一
		if len(permissions) != 0 && len(role.Permissions) != 0 && !containsAny(role.Permissions, permissions) {
			continue
		}

		// 拥有所有权限
		if dataScope == DATA_SCOPE_ALL {
			return db, nil
		}

		// 已经写入过数据
		if queryCondition.Len() > 0 {
			queryCondition.WriteString(" OR ")
		}

		if dataScope == DATA_SCOPE_CUSTOM {
			queryCondition.WriteString(fmt.Sprintf("%s.dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id = %d)", deptAlias, role.RoleId))
		} else if dataScope == DATA_SCOPE_DEPT {
			queryCondition.WriteString(fmt.Sprintf("%s.dept_id = %d", deptAlias, user.DeptId))
		} else if dataScope == DATA_SCOPE_DEPT_AND_CHILD {
			queryCondition.WriteString(fmt.Sprintf("%s.dept_id IN (SELECT dept_id FROM sys_dept WHERE dept_id = %d OR find_in_set( %d, ancestors ))", deptAlias, user.DeptId, user.DeptId))
		} else if dataScope == DATA_SCOPE_SELF {
			if userAlias != "" {
				queryCondition.WriteString(fmt.Sprintf("%s.user_id = %d", userAlias, user.UserId))
			} else {
				// 数据权限仅为本人并且没有userAlias别名不查询任何数据
				queryCondition.WriteString(fmt.Sprintf("%s.dept_id = 0", deptAlias))
			}
		}
		conditions.Add(dataScope)
	}

	// 多角色情况下所有角色都不包含传递过来的权限字符，此时不查询任何数据
	if conditions.Empty() {
		queryCondition.WriteString(fmt.Sprintf("%s.dept_id = 0", deptAlias))
	}

	return db.Where(queryCondition.String()), nil
}

// 判断collection 是否包含array中的任一元素
func containsAny(collection []string, array []string) bool {
	if len(collection) == 0 || len(array) == 0 {
		return false
	}
	for _, str := range array {
		if strutil.ContainsAny(str, collection...) {
			return true
		}
	}
	return false
}
