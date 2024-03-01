package cache

import (
	"sync"
)

// singleflight模式防止缓存雪崩
type SingleFlight interface {
	// key 表示缓存key
	// fn 表示执行的缓存动作
	Do(key string, fn func() (any, error)) (any, error)
	// 在Do的基础上返回是否共享缓存
	DoEx(key string, fn func() (any, error)) (any, bool, error)
}

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type flightGroup struct {
	// key 映射到调用
	calls map[string]*call
	lock  sync.Mutex
}

func NewSingleFlight() SingleFlight {
	return &flightGroup{calls: map[string]*call{}}
}

// 创建call调用
func (g *flightGroup) createCall(key string) (*call, bool) {
	g.lock.Lock()
	// 如果已经存在调用
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c := &call{}
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()
	return c, false
}

// 执行目标调用
func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()
	c.val, c.err = fn()
}

func (g *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	c, done := g.createCall(key)
	if done {
		return c.val, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *flightGroup) DoEx(key string, fn func() (any, error)) (any, bool, error) {
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, true, c.err
}
