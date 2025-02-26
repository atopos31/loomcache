package singleflight

import (
	"log"
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

type CallFunc = func() (any, error)

func (g *Group) Do(key string, fn CallFunc) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		log.Printf("cache hit singleflight, key: %v", key)
		return c.val, c.err
	}
	log.Printf("cache miss singleflight, key: %v", key)
	
	c := new(call)
	c.wg.Add(1)

	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
