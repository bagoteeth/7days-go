package geecache

import (
	"fmt"
	"log"
	"sync"
)

//缓存不存在时的回调，从数据源获取数据，添加到缓存
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (r GetterFunc) Get(key string) ([]byte, error) {
	return r(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g
}

func (r *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := r.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return r.load(key)
}

//不同场景可能会getFromPeer
func (r *Group) load(key string) (ByteView, error) {
	return r.getLocally(key)
}

func (r *Group) getLocally(key string) (ByteView, error) {
	//未命中则调用回调
	bytes, err := r.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	//回调结果添加至mainCache
	r.populateCache(key, value)
	return value, nil
}

func (r *Group) populateCache(key string, value ByteView) {
	r.mainCache.add(key, value)
}
