package main

import (
	. "./goose"
	"fmt"
	"log"
)

func main() {
	engine := New()

	engine.GET("/", func(ctx *Context) {
		ctx.Send("hello main home", 200)
	})
	engine.GET("/hello", func(ctx *Context) {
		ctx.Html("<h1>looks bigger ,right?</h1>", 200)
	})
	engine.GET("/user", func(ctx *Context) {
		var obj RawMap
		obj = make(RawMap)
		obj["username"] = ctx.Query("username")
		obj["msg"] = "successfully received!"
		ctx.Json(obj, 200)
	})
	engine.POST("/login", func(ctx *Context) {
		var obj RawMap
		obj = make(RawMap)
		obj["username"] = ctx.PostForm("username")
		obj["password"] = ctx.PostForm("password")
		obj["msg"] = "successfully received!"
		ctx.Json(obj, 200)
	})
	// test trie
	engine.GET("/hello/home", func(ctx *Context) {})
	engine.GET("/hello/:name/info", func(ctx *Context) {
		ctx.Html(fmt.Sprintf("<h2>url param: name=%s</h2>", ctx.Param("name")), 200)
	})
	engine.GET("/hello/home/wkk", func(ctx *Context) {
		ctx.Html(fmt.Sprintf("<h2>GET param: %s</h2>", ctx.Query("name")), 200)
	})
	engine.GET("/hello/home/wkk/ljq", func(ctx *Context) {
		ctx.Html(ctx.Path, 200)
	})
	engine.GET("/hello/home/wkk/baby", func(ctx *Context) {})
	engine.GET("/hello/home/ljq", func(ctx *Context) {})
	engine.GET("/about", func(ctx *Context) {})
	engine.GET("/hi/good", func(ctx *Context) {})
	engine.GET("/hi/bad", func(ctx *Context) {})

	// test group
	v1 := engine.Group("/v1")
	{
		v1.GET("/", func(ctx *Context) {
			ctx.Html("<h1>Hello Group Router</h1>", 200)
		})
		v1.GET("/hello", func(ctx *Context) {
			ctx.Send(fmt.Sprintf("hello %s, you're at %s\n", ctx.Query("name"), ctx.Path), 200)
		})
		// test nested group
		v2 := v1.Group("/nested/v2")
		v2.GET("/", func(ctx *Context) {
			ctx.Html("<h1>Hello Nested Group Router</h1>", 200)
		})
		v2.GET("/hello", func(ctx *Context) {
			ctx.Send(fmt.Sprintf("hello %s, you're at %s\n", ctx.Query("name"), ctx.Path), 200)
		})
	}
	log.Fatal(engine.BoostEngine("localhost:9999"))
}
