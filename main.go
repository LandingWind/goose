package main

import (
	. "./goose"
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

	log.Fatal(engine.BoostEngine(":9999"))
}
