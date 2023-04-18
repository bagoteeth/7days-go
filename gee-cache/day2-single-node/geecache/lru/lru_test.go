package lru

import (
	"testing"
	"time"
)

type String string

//value需要实现Len方法，获取占用内存大小
func (r String) Len() int {
	return len(r)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Upsert("k1", String("1234"))
	if v, ok := lru.Get("k1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache get k1 failed\n")
	}
	if _, ok := lru.Get("k2"); ok {
		t.Fatalf("cache miss k2 failed\n")
	}
}

func TestGC(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1+k2+k3+v1+v2+v3) - 1
	lru := New(int64(cap), func(key string, value Value) {
		t.Logf("delete key: %s value: %+v\n", key, value)
	})
	//先删k1，5s后删除k2, k3
	lru.Upsert(k1, String(v1))
	lru.Upsert(k2, String(v2))
	lru.Upsert(k3, String(v3))
	time.Sleep(5 * time.Second)
	lru.Upsert("k4", String("ababab"))
}
