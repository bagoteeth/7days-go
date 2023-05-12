package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	usedBytes int64
	//双向链表
	ll *list.List
	//value为双向链表的节点
	cache map[string]*list.Element
	//删除时的回调
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		usedBytes: 0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (r *Cache) Get(key string) (val Value, ok bool) {
	if ele, ok := r.cache[key]; ok {
		//最近访问移至队尾（约定front为队尾）
		r.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (r *Cache) RemoveOldest() {
	//获取队首元素（约定back为队首）
	ele := r.ll.Back()
	if ele != nil {
		r.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//删除cache映射关系
		delete(r.cache, kv.key)
		//更新总共占用内存
		r.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//回调
		if r.OnEvicted != nil {
			r.OnEvicted(kv.key, kv.value)
		}
	}
}

func (r *Cache) Upsert(key string, value Value) {
	if ele, ok := r.cache[key]; ok {
		r.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		r.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//链表元素包含key，相当于既可以从cache的key访问ll的ele，又可以从ll的ele访问cache的key
		ele := r.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		r.cache[key] = ele
		r.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for r.maxBytes != 0 && r.maxBytes < r.usedBytes {
		r.RemoveOldest()
	}
}

func (r *Cache) Len() int {
	return r.ll.Len()
}
