package goose

import (
	"fmt"
	"net/http"
)

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

func (router *Router) handle(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := router.handlerFuncMap[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
