package store

var client Factory

type Factory interface {
	GenTables() GenTableStore
	GenTableColumns() GenTableColumnStore
}

func Client() Factory {
	return client
}

func SetClient(c Factory) {
	client = c
}
