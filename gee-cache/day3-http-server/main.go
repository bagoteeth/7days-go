package main

import (
	"7days-go/gee-cache/day3-http-server/geecache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"bago":  "123",
	"teeth": "456",
	"mokou": "789",
}

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	addr := "localhost:9969"
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
