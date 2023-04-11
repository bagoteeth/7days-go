package gee

import "net/http"

type HandlerFunc func(ctx *Context)

//engine生成一个context，给router处理，依赖关系engine -> router -> context
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (r *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	r.router.addRoute(method, pattern, handler)
}

func (r *Engine) Get(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *Engine) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	r.router.handle(c)
}

func (r *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}
