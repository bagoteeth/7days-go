package gee

import (
	"fmt"
	"net/http"
)

type HanderFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HanderFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HanderFunc)}
}

func (r *Engine) addRoute(method string, pattern string, hander HanderFunc) {
	key := method + "-" + pattern
	r.router[key] = hander
}

func (r *Engine) GET(pattern string, hander HanderFunc) {
	r.addRoute("GET", pattern, hander)
}

func (r *Engine) POST(pattern string, hander HanderFunc) {
	r.addRoute("POST", pattern, hander)
}

func (r *Engine) RUN(addr string) error {
	return http.ListenAndServe(addr, r)
}

func (r *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	k := req.Method + "-" + req.URL.Path
	if h, ok := r.router[k]; ok {
		h(w, req)
	} else {
		fmt.Fprintf(w, "404 not found: %s\n", req.URL)
	}
}
