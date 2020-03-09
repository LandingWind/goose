package goose

import (
	"log"
	"net/http"
)

// context 封装req和res
type HandlerFunc func(*Context)

/*
 ** Engine接收器: 实现http.ListenAndServe中的handler接口
 */
type Engine struct {
	router *Router
}

/*
 ** func New() *Engine: 分配Engine的内存空间，返回指针供模块外使用
 */
func New() *Engine {
	log.SetPrefix("【GooseEngine】")
	return &Engine{router: newRouter()}
}

/*
 ** func BoostEngine(): 启动Engine, 即调用http.ListenAndServe(), 将默认handler设置为engine
 */
func (engine *Engine) BoostEngine(baseUrl string) (err error) {
	log.Println("Goose Engine Started...")
	return http.ListenAndServe(baseUrl, engine)
}

/*
 ** func GET(): 注册一个GET HandlerFunc
 */
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addHandlerFunc("GET", pattern, handler)
}

/*
 ** func POST(): 注册一个POST HandlerFunc
 */
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addHandlerFunc("POST", pattern, handler)
}

/*
 ** func ServeHTTP(): 接口函数ServeHTTP的具体实现
 ** 从map中检索是否注册有handler 调用handler或404 error
 */
func (engine *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	context := newContext(res, req)
	log.Println(req.Method + ": " + req.URL.String())
	engine.router.handle(context)
}
