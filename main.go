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
	log.Fatal(engine.BoostEngine(":9999"))
}
