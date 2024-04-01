package vo

type RouterVo struct {
	// 路由名字
	Name string `json:"name,omitempty"`
	// 路由地址
	Path string `json:"path,omitempty"`
	// 是否隐藏路由
	Hidden bool `json:"hidden,omitempty"`
	// 重定向地址，设置为noRedirect 时路由在面包屑导航中不可被点击
	Redirect string `json:"redirect,omitempty"`
	// 组件地址
	Component string `json:"component,omitempty"`
	// 路由参数，如{"id":1, "name":"sophie"}
	Query string `json:"query,omitempty"`
	// 给定路由下children 声明的路由大于1时，自动变成嵌套的模式
	AlwaysShow bool `json:"alwaysShow,omitempty"`
	// 其他元素
	Meta MetaVo `json:"meta,omitempty"`
	// 子路由
	Children []RouterVo `json:"children,omitempty"`
}
