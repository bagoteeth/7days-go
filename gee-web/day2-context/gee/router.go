package gee

import (
	"log"
	"net/http"
)

//负责根据传的context不同，进行不同的handle
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("route %s - %s\n", method, pattern)
	k := method + "-" + pattern
	r.handlers[k] = handler
}

func (r *router) handle(c *Context) {
	k := c.Method + "-" + c.Path
	if handler, ok := r.handlers[k]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
	}
}
