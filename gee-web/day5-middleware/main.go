package main

import (
	"7days-go/gee-web/day5-middleware/gee"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gee.New()
	r.Use(gee.Logger())
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
	v2.Use(V2Plug())
	{
		v2.Get("/hello", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "v2 hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
		v2.Get("/hello/:name", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "v2:name hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
	}

	r.Run(":9967")
}

func V2Plug() gee.HandlerFunc {
	return func(ctx *gee.Context) {
		t := time.Now()
		ctx.Fail(http.StatusInternalServerError, "Internal server error")
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
