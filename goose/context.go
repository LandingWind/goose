package goose

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	// 封装的req和res
	Res http.ResponseWriter
	Req *http.Request
	// 解析过后的key
	Method string // post,get一类
	Path   string // /hello
	// res info
	StatusCode int
}
type RawMap map[string]interface{} // string-anytype map

func newContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Res:        res,
		Req:        req,
		Method:     req.Method,
		Path:       req.URL.Path,
		StatusCode: 200,
	}
}

// get params :public
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// post params :public
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

// set res status code
func (ctx *Context) status(code int) {
	ctx.StatusCode = code
	ctx.Res.WriteHeader(code)
}

func (ctx *Context) GetHeaderValue(key string) (value string) {
	return ctx.Res.Header().Get(key)
}
func (ctx *Context) GetHeaderAll() http.Header {
	// ctx.Res.Header() will return the header map
	return ctx.Res.Header()
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
func (ctx *Context) Json(obj RawMap, code int) {
	ctx.header("Content-Type", "application/json")
	ctx.status(code)
	// NewEncoder defines the IO writer
	encoder := json.NewEncoder(ctx.Res)
	// encoder.Encode transfer rawType to JsonType
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Res, err.Error(), 500)
	}
}
