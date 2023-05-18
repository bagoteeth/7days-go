package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := Map{
		hash:     fn,
		replicas: replicas,
		keys:     nil,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return &m
}

func (r *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < r.replicas; i++ {
			hash := int(r.hash([]byte(strconv.Itoa(i) + key)))
			r.keys = append(r.keys, hash)
			r.hashMap[hash] = key
		}
	}
	sort.Ints(r.keys)
}

func (r *Map) Get(key string) string {
	if len(r.keys) == 0 {
		return ""
	}
	hash := int(r.hash([]byte(key)))
	idx := sort.Search(len(r.keys), func(i int) bool {
		return r.keys[i] >= hash
	})
	return r.hashMap[r.keys[idx%len(r.keys)]]
}
