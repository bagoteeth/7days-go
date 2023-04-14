package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(ctx *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	//所有group使用同一个engine instance
	engine *Engine
}

//engine生成一个context，给router处理，依赖关系engine -> router -> context
//最终会有1个engine，多个routergroup，engine存所有group，同时engine和group都可以加路由（engine继承group，虽然加路由是group的方法，
//但做事的是engine（所有group中唯一的engine实例）中的router）。
type Engine struct {
	*RouterGroup

	router *router

	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newGroup := &RouterGroup{
		prefix:      r.prefix + prefix,
		middlewares: nil,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (r *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := r.prefix + comp
	log.Printf("Route %s - %s\n", method, pattern)
	r.engine.router.addRoute(method, pattern, handler)
}

func (r *RouterGroup) Get(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *RouterGroup) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range r.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	r.router.handle(c)
}

func (r *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}
