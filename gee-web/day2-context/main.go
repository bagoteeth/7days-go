package main

import (
	"7days-go/gee-web/day2-context/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Get("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>hello teeth</h1>")
	})
	//curl http://127.0.0.1:9964/hello?name=bago
	r.Get("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})
	r.Get("/getall", func(ctx *gee.Context) {
		for k, v := range ctx.Req.Header {
			ctx.String(http.StatusOK, "header[%q] = %q\n", k, v)
		}
	})
	//curl http://127.0.0.1:9964/login -X POST -d "username=bago&password=teeth"
	r.POST("/login", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.Run(":9964")
}
