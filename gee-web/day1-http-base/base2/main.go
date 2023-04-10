package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
}

func (r *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Head[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 not found: %s\n", req.URL)
	}
}

//自定义handler，拦截所有http请求进行处理
func main() {
	eng := Engine{}
	log.Fatal(http.ListenAndServe(":9962", &eng))
}
