package vo

type MetaVo struct {
	// 设置路由在侧边栏和面包屑中显示的文字
	Title string `json:"title,omitempty"`
	// 设置路由的图标 (前端资源文件目录src/assets/icons/svg）
	Icon string `json:"icon,omitempty"`
	// 设置为true时不会被<keep-alive>缓存
	NoCache bool `json:"noCache,omitempty"`
	// 内链地址(http(s)://开头）
	Link string `json:"link,omitempty"`
}
