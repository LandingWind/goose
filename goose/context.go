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
	Method string            // post,get一类
	Path   string            // /hello
	Params map[string]string // 动态路由参数
	// res info
	StatusCode int
	// 中间件
	handlers      []HandlerFunc
	index         int
	middleStorage map[string]interface{} // 中间件可储存中间运算数据
}
type RawMap map[string]interface{} // string-anytype map

func newContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Res:           res,
		Req:           req,
		Method:        req.Method,
		Path:          req.URL.Path,
		StatusCode:    200,
		index:         -1,
		middleStorage: make(map[string]interface{}),
	}
}

// next()
func (ctx *Context) Next() { // 每个handler只能调用一次
	// 这里比较难以理解，必须用循环走完整条handler链
	for ctx.index++; ctx.index < len(ctx.handlers); ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}
func (ctx *Context) MiddleStore(key string, value interface{}) {
	ctx.middleStorage[key] = value
}
func (ctx *Context) GetMiddleStorage(key string) interface{} {
	return ctx.middleStorage[key]
}

// get route param :public
func (ctx *Context) Param(key string) string {
	val, ok := ctx.Params[key]
	if ok {
		return val
	}
	return ""
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

// fail
func (ctx *Context) Fail(err string, code int) {
	ctx.index = len(ctx.handlers)
	ctx.Json(RawMap{"message": err}, code)
}
