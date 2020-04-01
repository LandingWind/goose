package main

import (
	"fmt"
	. "goose/goose"
	"log"
	"net/http"
)

func main() {
	engine := New()

	engine.GET("/", func(ctx *Context) {
		ctx.Send("hello main home", 1000, 1002)
	})
	engine.GET("/hello", func(ctx *Context) {
		ctx.Html("<h1>looks bigger ,right?</h1>")
	})
	// test get query
	engine.GET("/user", func(ctx *Context) {
		var obj RawMap
		obj = make(RawMap)
		obj["username"] = ctx.Query("username")
		obj["msg"] = "successfully received!"
		ctx.Json(obj)
	})
	// test post postForm
	engine.POST("/login", func(ctx *Context) {
		var obj RawMap
		obj = make(RawMap)
		obj["username"] = ctx.PostForm("username")
		obj["password"] = ctx.PostForm("password")
		obj["msg"] = "successfully received!"
		ctx.Json(obj)
	})
	// test trie
	engine.GET("/hello/home", func(ctx *Context) {})
	engine.GET("/hello/:name/info", func(ctx *Context) {
		ctx.Html(fmt.Sprintf("<h2>url param: name=%s</h2>", ctx.Param("name")), 200)
	})
	engine.GET("/hello/home/wkk", func(ctx *Context) {
		ctx.Html(fmt.Sprintf("<h2>GET param: %s</h2>", ctx.Query("name")), 200)
	})
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
	// test middleware
	m1 := engine.Group("/performance")
	m1.Use(func(context *Context) {
		context.MiddleStore("info", "Kysoo is a handsome boy!")
	})
	m1.GET("/hello/:name", func(ctx *Context) {
		info := ctx.GetMiddleStorage("info")
		ctx.Html(fmt.Sprintf("<h2>Hello: %s, %s</h2>", ctx.Param("name"), info), 200)
	})
	// test recovery
	engine.GET("/err", func(ctx *Context) {
		arr := []string{"hello"}
		ctx.Send(arr[2], http.StatusOK)
	})

	// test engine options
	options := make(map[string]bool)
	options["logPrefix"] = true
	options["logRequest"] = true
	options["logRequestBody"] = false
	engine.SetOptions(options)

	// test static server
	engine.Static("assets", "testdata/assets")
	engine.Static("html", "testdata/html")

	// test template
	engine.LoadHTMLGlob("testdata/html/*")
	engine.GET("/render", func(ctx *Context) {
		ctx.HtmlTemplate("abc.html", RawMap{
			"title": "测试模版渲染",
		})
	})

	// boost engine
	log.Fatal(engine.BoostEngine("localhost:9999"))
}
