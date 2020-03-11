# goose

基于Golang+net/http的`原生` `极简易`Web框架

#### Basic Framework Struct
> 参考开源框架[Gin](https://github.com/gin-gonic/gin)和[Gee](https://github.com/geektutu/7days-golang)

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
- [ ] static file support
- [ ] multipart form request
- [ ] upload single file support 
- [ ] server log file support
- [ ] colored log

**version 0.4**

todo...

#### Usage Sample

basic
```golang
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

dynamic router
```golang
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

group router
```golang
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

use middleware
```golang
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