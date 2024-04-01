package store

var client Factory

type Factory interface {
	Transaction
	GenTables() GenTableStore
	GenTableColumns() GenTableColumnStore
}

type Transaction interface {
	Begin() Factory
	Commit() error
	Rollback() error
}

func Client() Factory {
	return client
}

func SetClient(c Factory) {
	client = c
}
