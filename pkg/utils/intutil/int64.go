package intutil

type Compare func(str int64, target int64) bool

func CompareAllInt64(compare Compare, target int64, searchint64s ...int64) bool {
	for i := range searchint64s {
		if !compare(searchint64s[i], target) {
			return false
		}
	}
	return true
}

func CompareAnyInt64(compare Compare, target int64, searchint64s ...int64) bool {
	for i := range searchint64s {
		if compare(searchint64s[i], target) {
			return true
		}
	}
	return false
}

func ContainsAnyInt64(str int64, searchint64s ...int64) bool {
	return CompareAnyInt64(func(str1, str2 int64) bool { return str1 == str2 }, str, searchint64s...)
}
