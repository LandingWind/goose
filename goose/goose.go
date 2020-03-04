package goose

import (
	"fmt"
	"net/http"
)

/*
 ** Engine接收器: 实现http.ListenAndServe中的handler接口
 ** type Handler interface {
		ServeHTTP(w ResponseWriter, r *Request)
	}
 ** 传递参数handlerFuncMap: 储存所有响应函数
*/
type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type Engine struct {
	handlerFuncMap map[string]HandlerFunc
}

/*
 ** func New() *Engine: 分配Engine的内存空间，返回指针供模块外使用
*/
func New() *Engine {
	return &Engine{handlerFuncMap: make(map[string]HandlerFunc)}
}

/*
 ** func BoostEngine(): 启动Engine, 即调用http.ListenAndServe(), 将默认handler设置为engine
*/
func (engine *Engine) BoostEngine(baseUrl string) (err error) {
	return http.ListenAndServe(baseUrl, engine)
} 

/*
 ** func ServeHTTP(): 接口函数ServeHTTP的具体实现
 ** 从map中检索是否注册有handler 调用handler或404 error
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.handlerFuncMap[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

/*
 ** func addHandlerFunc(): 注册HandlerFunc "method-pattern"作为key
*/
func (engine *Engine) addHandlerFunc(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.handlerFuncMap[key] = handler
}

/*
 ** func GET(): 注册一个GET HandlerFunc
*/
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addHandlerFunc("GET", pattern, handler)
}

/*
 ** func POST(): 注册一个POST HandlerFunc
*/
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addHandlerFunc("POST", pattern, handler)
}