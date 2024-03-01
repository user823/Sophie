package evict

import "container/list"

type lru struct {
	// 元素限制
	limit    int
	evicts   *list.List
	elements map[string]*list.Element
	// 淘汰回调函数
	onEvict func(key string)
}

type EmptyLRU struct{}

func (e EmptyLRU) Add(key string)    {}
func (e EmptyLRU) Remove(key string) {}
func (e EmptyLRU) Remains() int      { return 0 }
func (e EmptyLRU) Clean()            {}

func NewLRU(limit int, onEvict func(key string)) EvictPolicy {
	return &lru{
		limit:    limit,
		onEvict:  onEvict,
		evicts:   list.New(),
		elements: map[string]*list.Element{},
	}
}

func (l *lru) Add(key string) {
	if elem, ok := l.elements[key]; ok {
		l.evicts.MoveToFront(elem)
		return
	}

	elem := l.evicts.PushFront(key)
	l.elements[key] = elem

	if l.evicts.Len() > l.limit {
		l.removeOldest()
	}
}

func (l *lru) Remove(key string) {
	if elem, ok := l.elements[key]; ok {
		l.removeElement(elem)
	}
}

func (l *lru) removeElement(e *list.Element) {
	l.evicts.Remove(e)
	key := e.Value.(string)
	delete(l.elements, key)
	l.onEvict(key)
}

func (l *lru) removeOldest() {
	b := l.evicts.Back()
	if b != nil {
		l.removeElement(b)
	}
}

func (l *lru) Remains() int {
	return l.limit - l.evicts.Len()
}

func (l *lru) Clean() {
	l.evicts = list.New()
	l.elements = map[string]*list.Element{}
}
