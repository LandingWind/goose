package goose


type Router struct {
	handlerFuncMap map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{handlerFuncMap: make(map[string]HandlerFunc)}
}

/*
 ** func addHandlerFunc(): 注册HandlerFunc "method-pattern"作为key
*/
func (router *Router) addHandlerFunc(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlerFuncMap[key] = handler
}

func (router *Router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if handler, ok := router.handlerFuncMap[key]; ok {
		handler(ctx)
	} else {
		ctx.Send("Not Found", 404)
	}
}
