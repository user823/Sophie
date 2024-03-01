package v1

const (
	ServiceName = "Sophie Gen"
)

// 代码生成通用常量
const (
	// 单表（增删改查）
	TPL_CRUD = "crud"
	// 树表（增删改查）
	TPL_TREE = "tree"
	// 主子表（增删改查）
	TPL_SUB = "sub"
	// 树编码字段
	TREE_CODE = "treeCode"
	// 树父编码字段
	TREE_PARENT_CODE = "treeParentCode"
	// 树名称字段
	TREE_NAME = "treeName"
	// 上级菜单ID字段
	PARENT_MENU_ID = "parentMenuId"
	// 上级菜单名称字段
	PARENT_NENU_NAME = "parentMenuName"

	// 文本框
	HTML_INPUT = "input"
	// 文本域
	HTML_TEXTAREA = "textarea"
	// 下拉框
	HTML_SELECT = "select"
	// 单选框
	HTML_RADIO = "radio"
	// 复选框
	HTML_CHECKBOX = "checkbox"
	// 日期控件
	HTML_DATETIME = "datetime"
	// 图片上传控件
	HTML_IMAGE_UPLOAD = "imageUpload"
	// 文件上传控件
	HTML_FILE_UPLOAD = "fileUpload"
	// 富文本控件
	HTML_EDITOR = "editor"
	// 字符串类型
	TYPE_STRING = "string"
	// 整型
	TYPE_INT = "int"
	// 长整型
	TYPE_INT64 = "int64"
	// 浮点型
	TYPE_DOUBLE = "float64"
	// 时间类型
	TYPE_DATE = "Time"
	// 模糊查询
	QUERY_LIKE = "LIKE"
	// 相等查询
	QUERY_EQ = "EQ"
	// 需要
	REQUIRE = "1"
)

var (
	// 数据库字符串类型
	COLUMNTYPE_STR = []string{"char", "varchar", "nvarchar", "varchar2"}
	// 数据库文本类型
	COLUMNTYPE_TEXT = []string{"tinytext", "text", "mediumtext", "longtext"}
	// 数据库时间类型
	COLUMNTYPE_TIME = []string{"datetime", "time", "date", "timestamp"}
	// 数据库数字类型
	COLUMNTYPE_NUMBER = []string{"tinyint", "smallint", "mediumint", "int", "number", "integer", "bigint", "float", "double", "decimal"}
	// 页面不需要编辑字段
	COLUMNNAME_NOT_EDIT = []string{"id", "create_by", "create_time", "del_flag"}
	// 页面不需要显示的列表字段
	COLUMNNAME_NOT_LIST = []string{"id", "create_by", "create_time", "del_flag", "update_by", "update_time"}
	// 页面不需要查询字段
	COLUMNNAME_NOT_QUERY = []string{"id", "create_by", "create_time", "del_flag", "update_by",
		"update_time", "remark"}
	// Entity 基类字段
	BASE_ENTITY = []string{"CreateBy", "CreatedAt", "UpdateBy", "UpdatedAt", "Remark"}
	// Tree基类字段
	TREE_ENTITY = []string{"ParentName", "ParentId", "OrderNum", "Ancestors"}
)
