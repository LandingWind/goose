package goose

import (
	"net/http"
	// "encoding/json"
)


type Context struct {
	// 封装的req和res
	Res http.ResponseWriter
	Req *http.Request
	// 解析过后的key
	Method string // post,get一类
	Path string // /hello
	// res info
	StatusCode int
}

func newContext(res http.ResponseWriter, req *http.Request) *Context{
	return &Context{
		Res: res,
		Req: req,
		Method: req.Method,
		Path: req.URL.Path,
		StatusCode: 200,
	}
}

// get params :public
func (ctx *Context) Query() string {
	// return ctx.req.URL
	return ""
}
// post params :public
func (ctx *Context) Body() string {
	// return ctx.req
	return ""
}
// set res status code
func (ctx *Context) status(code int) {
	ctx.StatusCode = code
	ctx.Res.WriteHeader(code)
}

func (ctx *Context) header(key string, value string) {
	ctx.Res.Header().Set(key, value)
}

// plain text :public
func (ctx *Context) Send(txt string, code int) {
	ctx.header("Content-Type", "text/plain")
	ctx.status(code)
	ctx.Res.Write([]byte(txt))
}
// html :public
func (ctx *Context) Html(txt string, code int) {
	ctx.header("Content-Type", "text/html")
	ctx.status(code)
	ctx.Res.Write([]byte(txt))
}
// json
// func (ctx *Context) Json(obj interface{}, code int) {
// 	ctx.header("Content-Type", "application/json")
// 	ctx.status(code)
// 	encoder := json.NewEncoder(ctx.res)
// 	if err := encoder.Encode(obj); err != nil {
// 		http.Error(ctx,res, err.Error(), 500)
// 	}
// }