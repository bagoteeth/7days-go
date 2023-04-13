package main

import (
	"7days-go/gee-web/day4-group/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Get("/index", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>helloworld</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.Get("/", func(ctx *gee.Context) {
			ctx.HTML(http.StatusOK, "<h1>v1 root</h1>")
		})
		v1.Get("/hello", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "v1 hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.Get("/hello", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "v2 hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
		v2.Get("/hello/:name", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "v2:name hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
	}
	r.Run(":9966")
}
