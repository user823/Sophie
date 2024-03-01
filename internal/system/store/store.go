package store

var client Factory

type Factory interface {
	Transaction
	Users() UserStore
	UserPosts() UserPostStore
	UserRoles() UserRoleStore
	RoleMenus() RoleMenuStore
	Roles() RoleStore
	RoleDepts() RoleDeptStore
	Posts() PostStore
	OperLogs() OperLogStore
	Notices() NoticeStore
	Menus() MenuStore
	Logininfors() LogininforStore
	DictTypes() DictTypeStore
	DictData() DictDataStore
	Depts() DeptStore
	Configs() ConfigStore
	Close() error
}

type Transaction interface {
	Begin() Factory
	Commit() error
	Rollback() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
