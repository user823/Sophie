package vo

type RouterVo struct {
	// 路由名字
	Name string
	// 路由地址
	Path string
	// 是否隐藏路由
	Hidden bool
	// 重定向地址，设置为noRedirect 时路由在面包屑导航中不可被点击
	Redirect string
	// 组件地址
	Component string
	// 路由参数，如{"id":1, "name":"sophie"}
	Query string
	// 给定路由下children 声明的路由大于1时，自动变成嵌套的模式
	AlwaysShow bool
	// 其他元素
	Meta MetaVo
	// 子路由
	Children []RouterVo
}
