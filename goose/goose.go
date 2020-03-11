package goose

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	// 设置选项
	prefixString   string // 前缀字符串
	logPrefix      bool   // log是否有前缀
	logTime        bool   // log是否有时间显示
	logRouterTree  bool   // 启动时是否显示已注册handle的路径
	logRequest     bool   // log是否显示请求信息
	logRequestBody bool   // log是否显示请求体信息
	logPerformance bool   // log是否显示请求的响应时间
}

/*
 ** func New() *Engine: 分配Engine的内存空间，返回指针供模块外使用
 */
func New() *Engine {
	engine := &Engine{
		RouterGroup:    nil,
		router:         newRouter(),
		groups:         nil,
		prefixString:   "【Goose】",
		logPrefix:      false,
		logTime:        false,
		logRouterTree:  false,
		logRequest:     false,
		logRequestBody: false,
		logPerformance: false,
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
		// engine的RouterGroup.prefix为空字符串
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.Use(engine.Logger(), Recovery()) // 默认全局中间件
	log.SetFlags(0)
	return engine
}

func (engine *Engine) SetOptions(options interface{}) {
	val, ok := options.(map[string]bool)
	if !ok {
		fmt.Errorf("incorrect logOptions format")
		return
	}
	for key, item := range val {
		if key == "logPrefix" {
			engine.logPrefix = item
		} else if key == "logTime" {
			engine.logTime = item
		} else if key == "logRouterTree" {
			engine.logRouterTree = item
		} else if key == "logRequest" {
			engine.logRequest = item
		} else if key == "logRequestBody" {
			engine.logRequestBody = item
		} else if key == "logPerformance" {
			engine.logPerformance = item
		}
	}
	// 激活设置
	if engine.logPrefix {
		log.SetPrefix(engine.prefixString)
	} else {
		log.SetPrefix("")
	}
	if engine.logTime {
		log.SetFlags(log.LstdFlags)
	} else {
		log.SetFlags(0)
	}
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
	if group.engine.logRouterTree {
		log.Printf("Route %4s - %s", method, pattern)
	}
	group.engine.router.addHandlerFunc(method, pattern, handler)
}
func (group *RouterGroup) Use(handlers ...HandlerFunc) { // 可变参数添加middleware
	group.middleware = append(group.middleware, handlers...)
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
 ** func Any(): 注册一个响应所有method的HandlerFunc
 */
func (group *RouterGroup) Any(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
	group.addRoute("POST", pattern, handler)
	group.addRoute("PUT", pattern, handler)
	group.addRoute("DELETE", pattern, handler)
	group.addRoute("HEAD", pattern, handler)
	group.addRoute("PATCH", pattern, handler)
	group.addRoute("OPTIONS", pattern, handler)
}

/*
 ** func GET(): 注册一个GET HandlerFunc
 ** GET请求头没有Content-Type
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
 ** func DELETE(): 注册一个DELETE HandlerFunc，删除
 */
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

/*
 ** func PUT(): 注册一个PUT HandlerFunc，修改数据
 */
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}

/*
 ** func HEAD(): 注册一个HEAD HandlerFunc，获取报文首部
 */
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.addRoute("HEAD", pattern, handler)
}

/*
 ** func PATCH(): 注册一个PATCH HandlerFunc，不常用
 */
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.addRoute("PATCH", pattern, handler)
}

/*
 ** func OPTIONS(): 注册一个OPTIONS HandlerFunc，不常用
 */
func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.addRoute("OPTIONS", pattern, handler)
}

/*
 ** func ServeHTTP(): 接口函数ServeHTTP的具体实现
 ** 从map中检索是否注册有handler 调用handler或404 error
 */
func (engine *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	context := newContext(res, req) // pack context
	// get relevant middleware
	for _, item := range engine.groups {
		if strings.HasPrefix(context.Path, item.prefix) {
			context.handlers = append(context.handlers, item.middleware...)
		}
	}
	if engine.logRequest {
		log.Println(req.Method + ": " + req.URL.String())
	}
	engine.router.handle(context)
}
