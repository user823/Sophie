package db

// KeyValueStore 是所有k/v数据库存储后端的标准接口
type KeyValueStore interface {
	SetOptions(any)
	Connect() error
	Connected() bool
	Disconnect() error
	GetKey(string) (string, error)
	GetMultiKey([]string) ([]string, error)
	GetRawKey(string) (string, error)
	SetKey(string, string, int64) error
	SetRawKey(string, string, int64) error
	SetExp(string, int64) error
	GetExp(string) (int64, error)
	GetKeys(string) []string
	DeleteKey(string) bool
	DeleteAllKeys() bool
	DeleteRawKey(string) bool
	GetKeysAndValues() map[string]string
	GetKeysAndValuesWithFilter(string) map[string]string
	DeleteKeys([]string) bool
	Decrement(string)
	IncrememntWithExpire(string, int64) int64
	SetRollingWindow(key string, per int64, val string, pipeline bool) (int, []interface{})
	GetRollingWindow(key string, per int64, pipeline bool) (int, []interface{})
	GetSet(string) (map[string]string, error)
	AddToSet(string, string)
	GetAndDeleteSet(string) []interface{}
	RemoveFromSet(string, string)
	DeleteScanMatch(string) bool
	GetKeyPrefix() string
	AddToSortedSet(string, string, float64)
	GetSortedSetRange(string, string, string) ([]string, []float64, error)
	RemoveSortedSetRange(string, string, string) error
	GetListRange(string, int64, int64) ([]string, error)
	RemoveFromList(string, string) error
	AppendToSet(string, string)
	Exists(string) (bool, error)
}
