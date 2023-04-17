package main

import (
	"7days-go/gee-web/day5-middleware/gee"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.Use(Recovery())
	r.Get("/index", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>helloworld</h1>")
	})

	//数组越界
	r.Get("/panic", func(ctx *gee.Context) {
		names := []string{"bagoteeth"}
		ctx.String(http.StatusOK, names[100])
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

	r.Run(":9968")
}

func V2Plug() gee.HandlerFunc {
	return func(ctx *gee.Context) {
		t := time.Now()
		ctx.Fail(http.StatusInternalServerError, "Internal server error")
		log.Printf("[%d] %s in %v for group v2", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}

func Recovery() gee.HandlerFunc {
	return func(ctx *gee.Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(msg))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		//Recovery之后执行的handler如果panic，会被捕获
		ctx.Next()
	}
}

func trace(msg string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(msg + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
