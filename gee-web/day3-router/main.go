package main

import (
	"7days-go/gee-web/day3-router/gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.Get("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>helloworld</h1>")
	})
	r.Get("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})
	r.Get("/hello/:name", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})
	r.Get("/assets/*filepath", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{
			"filepath": ctx.Param("filepath"),
		})
	})

	r.Run(":9965")
}
