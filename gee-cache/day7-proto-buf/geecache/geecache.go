package geecache

import (
	"7days-go/gee-cache/day7-proto-buf/geecache/geecachepb"
	"7days-go/gee-cache/day7-proto-buf/geecache/singleflight"
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

	peers PeerPicker

	loader *singleflight.Group
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
		loader:    &singleflight.Group{},
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

func (r *Group) RegisterPeers(peers PeerPicker) {
	if r.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	r.peers = peers
}

//不同场景可能会getFromPeer
func (r *Group) load(key string) (value ByteView, err error) {
	viewi, err := r.loader.Do(key, func() (interface{}, error) {
		if r.peers != nil {
			if peer, ok := r.peers.PickPeer(key); ok {
				if value, err = r.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return r.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (r *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &geecachepb.Request{
		Group: r.name,
		Key:   key,
	}
	res := &geecachepb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
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
