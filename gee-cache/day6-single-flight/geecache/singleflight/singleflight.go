package singleflight

import "sync"

type call struct {
	sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	sync.Mutex
	m map[string]*call
}

func (r *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	r.Lock()
	if r.m == nil {
		r.m = make(map[string]*call)
	}
	if c, ok := r.m[key]; ok {
		r.Unlock()
		c.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.Add(1)
	r.m[key] = c
	r.Unlock()

	c.val, c.err = fn()
	c.Done()

	r.Lock()
	delete(r.m, key)
	r.Unlock()

	return c.val, c.err
}
