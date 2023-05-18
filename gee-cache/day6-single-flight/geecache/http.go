package geecache

import (
	"7days-go/gee-cache/day6-single-flight/geecache/consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self string
	//默认前缀
	basePath string

	sync.Mutex
	//一致性hash，更具key选择节点
	peers *consistenthash.Map
	//每个节点一个getter
	httpGetters map[string]*httpGetter
}

//每个server有多个，相当于client，
type httpGetter struct {
	baseURL string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (r *HTTPPool) Set(peers ...string) {
	r.Lock()
	defer r.Unlock()
	r.peers = consistenthash.New(defaultReplicas, nil)
	r.peers.Add(peers...)
	r.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		r.httpGetters[peer] = &httpGetter{baseURL: peer + r.basePath}
	}
}

func (r *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	r.Lock()
	defer r.Unlock()
	if peer := r.peers.Get(key); peer != "" && peer != r.self {
		r.Log("Pick peer %s", peer)
		return r.httpGetters[peer], true
	}
	return nil, false
}

func (r *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", r.self, fmt.Sprintf(format, v...))
}

func (r *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, r.basePath) {
		panic("HTTPPool unexpected path: " + req.URL.Path)
	}
	r.Log("%s %s", req.Method, r.basePath)
	parts := strings.SplitN(req.URL.Path[len(r.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

func (r *httpGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", r.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}
