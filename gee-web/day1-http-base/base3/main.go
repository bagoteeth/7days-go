package main

import (
	"7days-go/gee-web/day1-http-base/base3/gee"
	"fmt"
	"net/http"
)

//自定义client，集成了handler，支持对不同的url添加处理，http服务由client控制
func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "url.path = %q\n", r.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "header[%q] = %q\n", k, v)
		}
	})
	r.RUN(":9963")
}
