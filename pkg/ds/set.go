package ds

// 标记集合
type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return map[T]struct{}{}
}

func (s Set[T]) Add(keys ...T) {
	for i := range keys {
		s[keys[i]] = struct{}{}
	}
}

func (s Set[T]) Contains(key T) bool {
	if _, ok := s[key]; ok {
		return true
	}
	return false
}

func (s Set[T]) Remove(keys ...T) {
	for i := range keys {
		delete(s, keys[i])
	}
}

func (s Set[T]) Empty() bool {
	return len(s) == 0
}

func (s Set[T]) Values() (res []T) {
	res = make([]T, len(s))
	for k := range s {
		res = append(res, k)
	}
	return res
}
