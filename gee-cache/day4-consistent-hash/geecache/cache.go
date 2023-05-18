package geecache

import (
	"7days-go/gee-cache/day4-consistent-hash/geecache/lru"
	"sync"
)

//初始化时只需要cacheBytes
type cache struct {
	sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (r *cache) add(k string, v ByteView) {
	r.Lock()
	defer r.Unlock()
	//懒汉式
	if r.lru == nil {
		r.lru = lru.New(r.cacheBytes, nil)
	}
	r.lru.Upsert(k, v)
}

func (r *cache) get(k string) (v ByteView, ok bool) {
	r.Lock()
	defer r.Unlock()
	if r.lru == nil {
		return
	}
	if v, ok := r.lru.Get(k); ok {
		return v.(ByteView), ok
	}
	return
}
