package gee

import (
	"net/http"
	"strings"
)

type router struct {
	//key为method POST GET..., value为路由的trie树根节点
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//如果有动态路由，就匹配到动态路由为止
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

//可能包含动态路由
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

//实际访问的路由，不会有动态路由
func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		//parts为路由表中可能带动态路由的路由
		//searchParts为实际访问的路由
		//searchParts ["hello", "bago", "abc"]
		//parts ["hello", ":name", "abc"]
		//params {"name": "bago"}
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		if handler, ok := r.handlers[key]; ok {
			handler(c)
		} else {
			c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
		}
	} else {
		c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
	}
}
