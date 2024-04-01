package engine

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	v1 "github.com/user823/Sophie/api/domain/gen/v1"
	"github.com/user823/Sophie/pkg/ds"
	"github.com/user823/Sophie/pkg/errors"
	"github.com/user823/Sophie/pkg/log"
	"github.com/user823/Sophie/pkg/utils"
	"github.com/user823/Sophie/pkg/utils/strutil"
	"io"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
)

type TmplContext = map[string]any

// 从TmplContext 中获取string 数据
// 默认值为空
func GetStrDataFromTmplContext(t TmplContext, key string, trans func(value string) string) string {
	if trans == nil {
		trans = func(value string) string { return value }
	}
	if a, ok := t[key]; ok {
		if str, ok := a.(string); ok {
			return trans(str)
		}
	}
	return ""
}

type TemplateEngine struct {
	locker sync.RWMutex
	// 根据id 查询数据上下文
	mps    map[int64]TmplContext
	parser *template.Template
	Properties
}

type Properties struct {
	// 引擎名称
	Name string
	// 模版搜索路径
	// 包装器搜索指定目录（包含子目录）下所有"*.template" 文件，使用模版路径进行注册
	// 默认搜索当前路径
	TemplatePaths []string
	// 注册成功的模板
	Templates ds.Set[string]
}

func newProperties() *Properties {
	return &Properties{
		Name:          "root",
		TemplatePaths: []string{"."},
		Templates:     ds.NewSet[string](),
	}
}

type WrapperOption func(o *Properties)

func WithName(name string) WrapperOption {
	return func(o *Properties) {
		o.Name = name
	}
}

func WithSearchPaths(paths ...string) WrapperOption {
	return func(o *Properties) {
		o.TemplatePaths = paths
	}
}

// 初始化模版引擎
func NewEngine(opts ...WrapperOption) (wrapper *TemplateEngine) {
	// 初始化属性集
	properties := newProperties()
	for i := range opts {
		opts[i](properties)
	}
	log.Infof(" 当前的属性集合为: %v", properties)
	wrapper = &TemplateEngine{
		mps:        map[int64]TmplContext{},
		Properties: *properties,
	}

	wrapper.parser = template.New(wrapper.Properties.Name)
	// 加载引擎功能函数
	loadTemplateFuncs(wrapper.parser)

	// 搜索模板并加载
	for _, dir := range properties.TemplatePaths {
		dir = strings.Trim(dir, " ,")
		paths, err := utils.SearchFiles(dir, ".template")
		if err != nil {
			log.Errorf("template engine init loading path error: %s", err.Error())
			continue
		}
		for _, path := range paths {
			// 读取模板数据并加载
			data, err := os.ReadFile(path)
			if err != nil {
				log.Errorf("template engine read template file error: %s", err.Error())
				continue
			}

			// 不要使用ParseFile，因为不能处理不同目录文件同名情况
			registerPath := getRegisterName(path)
			_, err = wrapper.parser.New(registerPath).Parse(utils.B2s(data))
			if err != nil {
				log.Errorf("template engine parse template %s error: %s", path, err.Error())
				continue
			}
			wrapper.Templates.Add(registerPath)
		}
	}

	// 打印注册成功的模板
	log.Infof("注册成功的模板: %s", wrapper.parser.DefinedTemplates())
	return
}

// 引擎单例模式
// 全局模板引擎
var (
	engine *TemplateEngine
	once   sync.Once
)

func GetTemplateEngineOr(opts ...WrapperOption) *TemplateEngine {
	once.Do(func() {
		engine = NewEngine(opts...)
	})
	return engine
}

// 注册数据上下文
func (e *TemplateEngine) ParseContext(table *v1.GenTable) {
	dataContext := make(TmplContext)
	dataContext["id"] = table.TableId
	dataContext["tplCategory"] = table.TplCategory
	dataContext["tableName"] = table.Tablename
	dataContext["functionName"] = table.FunctionName
	dataContext["ClassName"] = table.ClassName
	dataContext["className"] = strutil.Uncapitalize(table.ClassName)
	dataContext["moduleName"] = table.ModuleName
	dataContext["BusinessName"] = strutil.Capitalize(table.BusinessName)
	dataContext["businessName"] = table.BusinessName
	dataContext["basePackage"] = getPackagePrefix(table.PackageName)
	dataContext["packageName"] = table.PackageName
	dataContext["author"] = table.FunctionAuthor
	dataContext["datetime"] = utils.Time2Str(time.Now())
	// pkColumn 是原始类型
	dataContext["pkColumn"] = table.PkColumn
	dataContext["importList"] = getImportList(table)
	dataContext["permissionPrefix"] = getPermissionPrefix(table.ModuleName, table.BusinessName)
	dataContext["columns"] = table.Columns
	dataContext["table"] = table
	dataContext["dicts"] = getDicts(table)
	setMenuTemplContext(dataContext, table)
	if v1.TPL_TREE == table.TplCategory {
		setTreeTmplContext(dataContext, table)
	}

	if v1.TPL_SUB == table.TplCategory {
		setSubTmplContext(dataContext, table)
	}

	// 数据修正
	if dataContext["functionName"] == "" {
		dataContext["functionName"] = "【请填写功能名称】"
	}

	e.locker.Lock()
	defer e.locker.Unlock()
	e.mps[table.TableId] = dataContext
}

// 模板渲染
func (e *TemplateEngine) ExecTemplate(template string, id int64) (result string, err error) {
	e.locker.RLock()
	data, ok := e.mps[id]
	e.locker.RUnlock()
	if !ok {
		return "", errors.New("模板渲染数据不存在")
	}

	buffer := bytes.NewBufferString(result)
	err = e.parser.ExecuteTemplate(buffer, template, data)
	return buffer.String(), err
}

// 渲染模板
func (e *TemplateEngine) ExecTemplateW(template string, id int64, w io.Writer) error {
	e.locker.RLock()
	data, ok := e.mps[id]
	e.locker.RUnlock()
	if !ok {
		return errors.New("模板渲染数据不存在")
	}

	err := e.parser.ExecuteTemplate(w, template, data)
	return err
}

// 获取模板文件列表
func (e *TemplateEngine) GetTemplateList(tplCategory, tplWebType string) (res []string) {
	useWebType := "vue"
	if tplWebType == "element-plus" {
		useWebType = "vue/v3"
	}
	templates := []string{"domain/domain.go.template", "store/store.go.template", "service/service.go.template", "app/*.template", "thrift/thrift.template",
		"sql/sql.template", "js/api.js.template"}
	if tplCategory == v1.TPL_CRUD {
		templates = append(templates, useWebType+"/index.vue.template")
	} else if tplCategory == v1.TPL_TREE {
		templates = append(templates, useWebType+"/index-tree.vue.template")
	} else if tplCategory == v1.TPL_SUB {
		templates = append(templates, useWebType+"/index.vue.template")
		templates = append(templates, "domain/sub-domain.vue.template")
	}

	registed := e.Templates.Values()
	// 从注册模板中搜索
	for i := range templates {
		if strutil.CompareAny(strutil.SimpleMatch, templates[i], registed...) {
			res = append(res, templates[i])
		}
	}
	return
}

// 设置目录模板上下文
func setMenuTemplContext(dataContext TmplContext, table *v1.GenTable) {
	if table.Options == "" {
		return
	}
	var result TmplContext
	if err := jsoniter.Unmarshal(utils.S2b(table.Options), &result); err != nil {
		log.Errorf("set menu template context error: %s", err.Error())
		return
	}
	dataContext["parentMenuId"] = getParentMenuId(result)
}

// 获取上级菜单id
func getParentMenuId(dataContext TmplContext) string {
	if a, ok := dataContext[v1.PARENT_MENU_ID]; ok {
		if str, ok := a.(string); ok {
			return str
		}
	}
	return v1.DEFAULT_PARENT_MENU_ID
}

// 设置树表结构上下文
func setTreeTmplContext(dataContext TmplContext, table *v1.GenTable) {
	var result TmplContext
	if err := jsoniter.Unmarshal(utils.S2b(table.Options), &result); err != nil {
		log.Errorf("set tree template context error: %s", err.Error())
		return
	}
	// 获取树编码
	dataContext["treeCode"] = GetStrDataFromTmplContext(result, v1.TREE_CODE, strutil.ToCamelCase)
	// 获取树父编码
	dataContext["treeParentCode"] = GetStrDataFromTmplContext(result, v1.TREE_PARENT_CODE, strutil.ToCamelCase)
	// 获取树名
	dataContext["treeName"] = GetStrDataFromTmplContext(result, v1.TREE_NAME, strutil.ToCamelCase)
	// 获取需要在哪一列上面显示展开按钮
	num := 0
	for i := range table.Columns {
		if table.Columns[i].IsList == v1.COLUMN_LIST_TYPE {
			num++
			if table.Columns[i].ColumnName == GetStrDataFromTmplContext(result, v1.TREE_NAME, nil) {
				break
			}
		}
	}
	dataContext["expandColumn"] = num
	if _, ok := result[v1.TREE_PARENT_CODE]; ok {
		dataContext["tree_parent_code"] = GetStrDataFromTmplContext(result, v1.TREE_PARENT_CODE, nil)
	}
	if _, ok := result[v1.TREE_NAME]; ok {
		dataContext["tree_name"] = GetStrDataFromTmplContext(result, v1.TREE_NAME, nil)
	}
}

// 设置主子表上下文结构
func setSubTmplContext(dataContext TmplContext, table *v1.GenTable) {
	dataContext["subTable"] = table.SubTable
	dataContext["subTableName"] = table.SubTableName
	dataContext["subTableFkName"] = table.SubTableFkName
	dataContext["subTableFkclassName"] = strutil.Uncapitalize(table.SubTableFkName)
	dataContext["subClassName"] = table.SubTable.ClassName
	dataContext["subclassName"] = strutil.Uncapitalize(table.SubTable.ClassName)
	dataContext["subImportList"] = getImportList(table.SubTable)
}

// 获取包前缀名
func getPackagePrefix(packageName string) string {
	index := strings.LastIndex(packageName, ".")
	if index == -1 {
		return ""
	}
	return packageName[:index]
}

// 根据列类型的特殊类型获取导入包
// 目前除基本类型外就只有时间类型
func getImportList(table *v1.GenTable) (importList []string) {
	columns := table.Columns
	for i := range columns {
		if !isSuperColumn(columns[i].GoField) && v1.TYPE_DATE == columns[i].GoType {
			importList = append(importList, "time")
		}
	}
	return
}

// 根据列类型获取字典组
func getDicts(table *v1.GenTable) string {
	dicts := ds.NewSet[string]()
	addDicts(dicts, table.Columns)
	if table.SubTable != nil {
		addDicts(dicts, table.SubTable.Columns)
	}
	return strings.Join(dicts.Values(), ",")
}

// 添加字典列表 （每个元素用”围起来）
func addDicts(dicts ds.Set[string], columns []*v1.GenTableColumn) {
	for i := range columns {
		if !isSuperColumn(columns[i].GoField) && columns[i].DictType != "" && strutil.ContainsAny(columns[i].HtmlType, v1.HTML_SELECT, v1.HTML_RADIO, v1.HTML_CHECKBOX) {
			dicts.Add("'" + columns[i].DictType + "'")
		}
	}
}

// 获取权限前缀
func getPermissionPrefix(moduleName, businessName string) string {
	return fmt.Sprintf("%s:%s", moduleName, businessName)
}

// 特殊处理的列类型(避免生成多余的属性）
func isSuperColumn(goField string) bool {
	return strutil.ContainsAny(goField, "createBy", "createTime", "updateBy", "updateTime", "remark",
		"parentName", "parentId", "orderNum", "ancestors")
}

// 根据路径获取注册模板名
func getRegisterName(path string) string {
	rootPaths := []string{"app/", "domain/", "js/", "service/", "sql/", "store/", "thrift/", "vue/v3/", "vue/"}
	for _, p := range rootPaths {
		index := strings.LastIndex(path, p)
		if index != -1 {
			return path[index:]
		}
	}
	return path
}
