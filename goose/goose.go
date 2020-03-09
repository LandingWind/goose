package goose

import (
	"log"
	"net/http"
)

// context 封装req和res
type HandlerFunc func(*Context)
type RouterGroup struct {
	prefix     string
	engine     *Engine // 指向同一个实例
	middleware []HandlerFunc
	parent     *RouterGroup // group嵌套
}

/*
 ** Engine接收器: 实现http.ListenAndServe中的handler接口
 */
type Engine struct {
	*RouterGroup // 匿名字段
	router       *Router
	groups       []*RouterGroup // 储存所有的group，包括嵌套的group
}

/*
 ** func New() *Engine: 分配Engine的内存空间，返回指针供模块外使用
 */
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
		// engine的RouterGroup.prefix为空字符串
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	log.SetPrefix("【Goose】")
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup) // 无论嵌套全部放在engine下
	return newGroup
}
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addHandlerFunc(method, pattern, handler)
}

/*
 ** func BoostEngine(): 启动Engine, 即调用http.ListenAndServe(), 将默认handler设置为engine
 */
func (engine *Engine) BoostEngine(baseUrl string) (err error) {
	log.Println("Goose Engine Started...")
	//log.Println("Show Trie Tree Routes:")
	//engine.router.printRouter()
	return http.ListenAndServe(baseUrl, engine)
}

/*
 ** func GET(): 注册一个GET HandlerFunc
 */
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

/*
 ** func POST(): 注册一个POST HandlerFunc
 */
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
