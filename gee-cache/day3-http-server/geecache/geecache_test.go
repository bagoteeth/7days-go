package geecache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"bago":  "123",
	"teeth": "456",
	"mokou": "789",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB search key]", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return cloneBytes([]byte(v)), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	//把所有数据写到cache
	for k, v := range db {
		//从缓存和数据源获取数据，不应该报错
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		//loadCounts > 1说明没有在第一次get就把数据源写到缓存
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
