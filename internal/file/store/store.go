package store

var client Factory

type Factory interface {
	Files() FileStore
}

func Client() Factory {
	return client
}

func SetClient(c Factory) {
	client = c
}
