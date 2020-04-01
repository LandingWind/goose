# goose

Based on Golang+net/http的 `native` `handiest` Web Framework

#### Basic Framework Struct
> Ideas parcially from [Gin](https://github.com/gin-gonic/gin) and [Gee](https://github.com/geektutu/7days-golang)

#### Feature

**version 0.1**

- [x] static router
- [x] dynamical router (based on TrieTree)
- [x] group router
- [x] error handler support
- [x] middleware support
- [x] middleware temporary data pass support
- [x] json,plain text,html support

**version 0.2 (mending version)**
- [x] more request method support
- [x] default value support
- [x] detailed format log support

**version 0.3**
- [x] static file support
- [x] template render support
- [ ] multipart form request
- [ ] upload single file support 
- [ ] server log file support
- [ ] colored log

**version 0.4**

todo...

#### Usage Sample

#### basic

```go
import (
	. "./goose"
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
  log.Fatal(engine.BoostEngine("localhost:9999"))
}
```

#### dynamic router

```go
import (
	. "./goose"
)
func main() {
	engine := New()
	engine.GET("/hello/home", func(ctx *Context) {})
  engine.GET("/hello/:name/info", func(ctx *Context) {
      ctx.Html(fmt.Sprintf("<h2>url param: name=%s</h2>", ctx.Param("name")), 200)
  })
  log.Fatal(engine.BoostEngine("localhost:9999"))
}
```

#### group router

```go
import (
	. "./goose"
)
func main() {
	engine := New()
	v1 := engine.Group("/v1")
  {
    v1.GET("/", func(ctx *Context) {
      ctx.Html("<h1>Hello Group Router</h1>", 200)
    })
    v1.GET("/hello", func(ctx *Context) {
      ctx.Send(fmt.Sprintf("hello %s, you're at %s\n", ctx.Query("name"), ctx.Path), 200)
    })
    // nested group
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
```

#### use middleware

```go
import (
	. "./goose"
)
func main() {
	engine := New()
	m1 := engine.Group("/performance")
  m1.Use(func(context *Context) {
    context.MiddleStore("info", "Kysoo is a handsome boy!")
  })
  m1.GET("/hello/:name", func(ctx *Context) {
    info := ctx.GetMiddleStorage("info")
    ctx.Html(fmt.Sprintf("<h2>Hello: %s, %s</h2>", ctx.Param("name"), info), 200)
  })
  log.Fatal(engine.BoostEngine("localhost:9999"))
}
```

#### set goose engine options

```go
/*
* 设置选项
prefixString   string // 前缀字符串
logPrefix      bool   // log是否有前缀
logTime        bool   // log是否有时间显示
logRouterTree  bool   // 启动时是否显示已注册handle的路径
logRequest     bool   // log是否显示请求信息
logRequestBody bool   // log是否显示请求体信息
logPerformance bool   // log是否显示请求的响应时间
*/
options := make(map[string]bool)
options["logPrefix"] = true
options["logRequest"] = true
engine.SetOptions(options)
```

#### set static file server

```go
// when visit ":port/assets/whatever file" 
// it will be reflected to ":port/testdata/assets/whatever file" 
// and process as static file 
engine.Static("assets","testdata/assets")
// btw, you can set multiple relfecting static path
```

#### html template

```go
// example code
engine.LoadHTMLGlob("testdata/html/*") // for loading templates into memory
engine.GET("/render", func(ctx *Context) {
    ctx.HtmlTemplate("abc.html", RawMap{
        "title": "测试模版渲染",
    })
})
// also, you can set renderer mapfunc by following way
r.SetFuncMap(template.FuncMap{
    "func1": func1,
    "func2": func2,
})
```