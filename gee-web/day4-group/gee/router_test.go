package gee

import (
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	t.Log(parsePattern("/p/:name"))
	t.Log(parsePattern("/p/:name/abc"))
	t.Log(parsePattern("/p/:name/:abc"))
	t.Log(parsePattern("/p/:name/*abc"))
	t.Log(parsePattern("/p/:name/:abc/dd"))
	t.Log(parsePattern("/p/:name/*abc/dd"))
	t.Log(parsePattern("/p/*"))
	t.Log(parsePattern("/p/*name"))
	t.Log(parsePattern("/p/*name/abc"))
	t.Log(parsePattern("/p/*name/:abc"))
	t.Log(parsePattern("/p/*name/*abc"))
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, p := r.getRoute("GET", "/hello/bago")
	if n == nil {
		t.Fatal("n nil")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("n pattern error")
	}
	if p["name"] != "bago" {
		t.Fatal("n p error: " + p["name"])
	}
	t.Logf("match path: %s, params['name'] = %s\n", n.pattern, p["name"])
}
