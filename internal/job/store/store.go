package store

var client Factory

type Factory interface {
	Jobs() JobStore
	JobLogs() JobLogStore
}

func Client() Factory {
	return client
}

func SetClient(c Factory) {
	client = c
}
