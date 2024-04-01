package engine

import (
	flag "github.com/spf13/pflag"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"strconv"
	"strings"
)

type GenHelperOptions struct {
	Author        string   `json:"author" mapstructure:"author"`
	PackageName   string   `json:"package_name" mapstructure:"package_name"`
	AutoRemovePre bool     `json:"auto_remove_pre" mapstructure:"auto_remove_pre"`
	TablePrefix   []string `json:"table_prefix" mapstructure:"table_prefix"`
}

func NewGenHelperOptions() *GenHelperOptions {
	return &GenHelperOptions{
		AutoRemovePre: false,
		TablePrefix:   []string{"sys_"},
	}
}

func (o *GenHelperOptions) Validate() error { return nil }

func (o *GenHelperOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Author, "gen.params.author", o.Author, ""+
		"Code Generator set default author")
	fs.StringVar(&o.PackageName, "gen.params.package_name", o.PackageName, ""+
		"Code Generator set default package_name")
	fs.BoolVar(&o.AutoRemovePre, "gen.params.auto_remove_pre", o.AutoRemovePre, ""+
		"Code Generator enable remove table prefix.")
	fs.StringSliceVar(&o.TablePrefix, "gen.params.table_prefix", o.TablePrefix, ""+
		"Code Generator search table prefixes")
}

// 全局配置
var (
	GlobalOptions GenHelperOptions
)

type GenHelper struct {
	config *GenHelperOptions
}

// 通过全局配置生成GenHelper
func DefaultGenHelper() *GenHelper {
	return &GenHelper{&GlobalOptions}
}

func NewGenHelper(config *GenHelperOptions) *GenHelper {
	return &GenHelper{config: config}
}

// 初始化表信息（用于导入）
func (g *GenHelper) InitTable(table *v1.GenTable, operName string) {
	table.ClassName = g.convertClassName(table.Tablename)
	table.PackageName = g.getPackageName(table)
	table.ModuleName = g.getModuleName(table)
	table.BusinessName = g.getBusinessName(table)
	table.FunctionName = table.TableComment
	table.FunctionAuthor = g.getAuthor(table)
	table.CreateBy = operName
}

// 初始化列属性字段
func (g *GenHelper) InitColumnField(table *v1.GenTable, column *v1.GenTableColumn) {
	dataType := g.getDbType(column)
	column.TableId = table.TableId
	column.CreateBy = table.CreateBy
	// 设置go 字段名
	column.GoField = strutil.ConvertToCamelCase(column.ColumnName)
	// 设置默认类型
	column.GoType = v1.TYPE_STRING
	column.QueryType = v1.QUERY_EQ

	if strutil.ContainsAny(dataType, v1.COLUMNTYPE_STR...) {
		// 字符串长度超过500则设置为文本域
		length := g.getColumnTextLength(dataType)
		if length >= 500 {
			column.HtmlType = v1.HTML_TEXTAREA
		} else {
			column.HtmlType = v1.HTML_INPUT
		}
	} else if strutil.ContainsAny(dataType, v1.COLUMNTYPE_TEXT...) {
		column.HtmlType = v1.HTML_TEXTAREA
	} else if strutil.ContainsAny(dataType, v1.COLUMNTYPE_TIME...) {
		column.GoType = v1.TYPE_DATE
		column.HtmlType = v1.HTML_DATETIME
	} else if strutil.ContainsAny(dataType, v1.COLUMNTYPE_NUMBER...) {
		column.HtmlType = v1.HTML_INPUT
		column.GoType = v1.TYPE_INT
	} else if strutil.ContainsAny(dataType, v1.COLUMNTYPE_FLOAT...) {
		column.HtmlType = v1.HTML_INPUT
		column.GoType = v1.TYPE_DOUBLE
	}

	// 插入字段 (默认所有字段都需要插入）
	column.IsInsert = v1.REQUIRE

	// 编辑字段
	if !strutil.ContainsAny(column.ColumnName, v1.COLUMNNAME_NOT_EDIT...) && !(column.IsPk != "" && column.IsPk != "1") {
		column.IsEdit = v1.REQUIRE
	}

	// 列表字段
	if !strutil.ContainsAny(column.ColumnName, v1.COLUMNNAME_NOT_LIST...) && !(column.IsPk != "" && column.IsPk != "1") {
		column.IsList = v1.REQUIRE
	}

	// 查询字段
	if !strutil.ContainsAny(column.ColumnName, v1.COLUMNNAME_NOT_QUERY...) && !(column.IsPk != "" && column.IsPk != "1") {
		column.IsQuery = v1.REQUIRE
	}

	// 查询字段类型
	if strutil.EndsWithIgnoreCase(column.ColumnName, "name") {
		column.QueryType = v1.QUERY_LIKE
	}

	// 状态字段设置单选框
	if strutil.EndsWithIgnoreCase(column.ColumnName, "status") {
		column.HtmlType = v1.HTML_RADIO
	} else if strutil.EndsWithIgnoreCase(column.ColumnName, "type") || strutil.EndsWithIgnoreCase(column.ColumnName, "sex") {
		// 类型&性别字段设置下拉框
		column.HtmlType = v1.HTML_SELECT
	} else if strutil.EndsWithIgnoreCase(column.ColumnName, "image") {
		column.HtmlType = v1.HTML_IMAGE_UPLOAD
	} else if strutil.EndsWithIgnoreCase(column.ColumnName, "file") {
		column.HtmlType = v1.HTML_FILE_UPLOAD
	} else if strutil.EndsWithIgnoreCase(column.ColumnName, "content") {
		column.HtmlType = v1.HTML_EDITOR
	}
}

// 根据表明转化为类名
func (g *GenHelper) convertClassName(tableName string) string {
	if g.config.AutoRemovePre && len(g.config.TablePrefix) > 0 {
		// 去除tableName 匹配到的第一个表前缀
		for _, prefix := range g.config.TablePrefix {
			if strings.HasPrefix(tableName, prefix) {
				tableName = tableName[len(prefix):]
				break
			}
		}
	}
	return strutil.ConvertToCamelCase(tableName)
}

// 获取包名
func (g *GenHelper) getPackageName(table *v1.GenTable) string {
	if table.PackageName == "" {
		return g.config.PackageName
	}
	return table.PackageName
}

// 获取模块名
func (g *GenHelper) getModuleName(table *v1.GenTable) string {
	index := strings.LastIndex(table.PackageName, ".")
	return table.PackageName[index+1:]
}

// 获取业务名
func (g *GenHelper) getBusinessName(table *v1.GenTable) string {
	if table.BusinessName != "" {
		return table.BusinessName
	}
	index := strings.LastIndex(table.Tablename, "_")
	return table.Tablename[index+1:]
}

func (g *GenHelper) getAuthor(table *v1.GenTable) string {
	if table.FunctionAuthor == "" {
		return g.config.Author
	}
	return table.FunctionAuthor
}

// 获取数据库类型字段
func (g *GenHelper) getDbType(column *v1.GenTableColumn) string {
	if strings.Contains(column.ColumnType, "(") {
		index := strings.Index(column.ColumnType, "(")
		return column.ColumnType[:index]
	}
	return column.ColumnType
}

// 获取文本字段长度
func (g *GenHelper) getColumnTextLength(columnType string) int {
	if strings.Contains(columnType, "(") {
		index := strings.Index(columnType, "(")
		parse, err := strconv.Atoi(columnType[index+1 : len(columnType)-1])
		if err != nil {
			log.Warnf("")
			return 0
		}
		return parse
	}
	return 0
}
