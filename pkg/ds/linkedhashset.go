package ds

import "container/list"

// T 为key的类型
type LinkedHashSet[T comparable] struct {
	mp   map[T]*list.Element
	list *list.List
}

type pair struct {
	key   any
	value any
}

func (l *LinkedHashSet[T]) Size() int {
	return len(l.mp)
}

func NewLinkedHashSet[T comparable]() *LinkedHashSet[T] {
	return &LinkedHashSet[T]{
		mp:   map[T]*list.Element{},
		list: list.New(),
	}
}

func (l *LinkedHashSet[T]) Empty() bool {
	return len(l.mp) == 0
}

func (l *LinkedHashSet[T]) Contains(key T) bool {
	if _, ok := l.mp[key]; ok {
		return true
	}
	return false
}

func (l *LinkedHashSet[T]) Find(key T) any {
	if a, ok := l.mp[key]; ok {
		return a.Value.(pair).value
	}
	return nil
}

func (l *LinkedHashSet[T]) Delete(key T) {
	if a, ok := l.mp[key]; ok {
		l.list.Remove(a)
		delete(l.mp, key)
	}
}

func (l *LinkedHashSet[T]) Front() any {
	return l.list.Front().Value.(pair).value
}

func (l *LinkedHashSet[T]) Back() any {
	return l.list.Back().Value.(pair).value
}

func (l *LinkedHashSet[T]) Add(key T, value any) {
	// 元素不存在时添加
	if _, ok := l.mp[key]; !ok {
		e := l.list.PushBack(pair{key, value})
		l.mp[key] = e
	}
}

// 移除队首元素
func (l *LinkedHashSet[T]) Poll() any {
	if len(l.mp) > 0 {
		e := l.list.Front()
		l.list.Remove(e)
		delete(l.mp, e.Value.(pair).key.(T))
		return e.Value.(pair).value
	}
	return nil
}

// 获取所有元素
func (l *LinkedHashSet[T]) AllItems() []any {
	result := make([]any, 0, len(l.mp))
	for _, v := range l.mp {
		result = append(result, v.Value.(pair).value)
	}
	return result
}
