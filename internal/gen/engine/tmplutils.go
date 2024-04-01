package engine

import (
	"fmt"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"path/filepath"
	"strings"
)

// 获取生成文件名
func GetFileName(template string, table *v1.GenTable) (fileName string) {
	goPath := filepath.Join(v1.PROJECT_ROOT, strings.ReplaceAll(table.PackageName, ".", string(filepath.Separator)))
	domainPath := filepath.Join(v1.DOMAIN_ROOT, table.ModuleName, v1.CODE_VERSION)
	cmdPath := v1.CMD_ROOT
	vuePath := v1.VUE_ROOT
	idlPath := filepath.Join(v1.IDL_ROOT, table.ModuleName, v1.CODE_VERSION)

	if strings.Contains(template, "sub-domain.go.template") {
		// 生成数据库交互实体类
		fileName = filepath.Join(domainPath, strings.ToLower(table.SubTable.ClassName)+".go")
	} else if strings.Contains(template, "domain.go.template") {
		// 生成数据库交互实体类
		fileName = filepath.Join(domainPath, strings.ToLower(table.ClassName)+".go")
	} else if strings.Contains(template, "sub-store.go.template") {
		// 生成存储层文件
		fileName = filepath.Join(goPath, "store/store.go")
	} else if strings.Contains(template, "store.go.template") {
		// 生成存储层文件
		fileName = filepath.Join(goPath, fmt.Sprintf("store/%s.go", GetStoreFileName(table.Tablename)))
	} else if strings.Contains(template, "service.go.template") {
		// 生成服务层文件
		fileName = filepath.Join(goPath, fmt.Sprintf("service/%s.go", table.BusinessName))
	} else if strings.Contains(template, "app") {
		baseName := filepath.Base(template)
		baseName = baseName[:strings.Index(baseName, ".")]
		if baseName == "main" {
			// 生成 app 启动类文件
			fileName = filepath.Join(cmdPath, baseName+".go")
		} else {
			// 生成app架构文件
			fileName = filepath.Join(goPath, baseName+".go")
		}
	} else if strings.Contains(template, "thrift.template") {
		// 生成微服务定义的idl文件
		fileName = filepath.Join(idlPath, table.ModuleName+".thrift")
	} else if strings.Contains(template, "sql.template") {
		// 生成菜单控制的sql
		fileName = table.BusinessName + "Menu.sql"
	} else if strings.Contains(template, "api.js.template") {
		// 生成前端请求api
		fileName = filepath.Join(vuePath, "api", table.ModuleName, table.BusinessName+".js")
	} else if strings.Contains(template, "index.vue.template") {
		// 生成index.vue
		fileName = filepath.Join(vuePath, "views", table.ModuleName, table.BusinessName, "index.vue")
	} else if strings.Contains(template, "index-tree.vue.template") {
		// 生成index-tree.vue
		fileName = filepath.Join(vuePath, "views", table.ModuleName, table.BusinessName, "index.vue")
	}
	return fileName
}

// 根据表名获取存储层文件名
func GetStoreFileName(tableName string) string {
	// 没有设置表名时返回默认名store
	if tableName == "" {
		return "store"
	}

	// 检查表名是否以"sys_"开头
	if strings.HasPrefix(tableName, v1.TABLENAME_PREFIX) {
		tableName = tableName[len(v1.TABLENAME_PREFIX):]
	}

	// 转化为驼峰表示法
	tableName = strutil.ToCamelCase(tableName)
	return strings.ToLower(tableName)
}
