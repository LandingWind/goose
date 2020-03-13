package goose

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	rootNode       map[string]*TrieNode
	handlerFuncMap map[string]HandlerFunc // golang函数一般是值调用
	assetsMap      map[string]string      // 记录静态文件路径
}

// get, post each has a Trie Router Tree

func newRouter() *Router {
	return &Router{
		handlerFuncMap: make(map[string]HandlerFunc),
		rootNode:       make(map[string]*TrieNode),
		assetsMap:      make(map[string]string),
	}
}

func (router *Router) setStatic(alias string, path string) {
	if alias == "" {
		router.assetsMap[path] = path
	} else {
		router.assetsMap[alias] = path
	}
}

/*
 ** func parsePattern(): 解析路径 返回短路径数组
 */
func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	urlSplit := strings.Split(pattern, "/") // return string[]
	for _, val := range urlSplit {
		if val != "" {
			parts = append(parts, val)
		}
	}
	return parts
}

/*
 ** func addHandlerFunc(): 注册HandlerFunc "method-pattern"作为key
 */
func (router *Router) addHandlerFunc(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	// 判断是否存在method对应的trie树
	_, ok := router.rootNode[method]
	if !ok {
		// 没有method根节点于是创建一个
		router.rootNode[method] = &TrieNode{}
	}
	// 将路径解析为parts
	parts := parsePattern(pattern)
	// insert path 更新trie树
	router.rootNode[method].insertPath(pattern, parts, 0)
	router.handlerFuncMap[key] = handler
}

func (router *Router) handle(ctx *Context) {
	searchParts := parsePattern(ctx.Path)
	flag := false
	if len(searchParts) > 0 {
		first := searchParts[0]
		if _, ok := router.assetsMap[first]; ok {
			flag = true
			//file := ctx.Path[len(first)+1:]
			//fmt.Println("val: " + val)
			//fmt.Println("file: " + file)
			fs := http.Dir(router.assetsMap[first])
			prefix := "/" + first + "/"
			fileServer := http.StripPrefix(prefix, http.FileServer(fs))
			ctx.handlers = append(ctx.handlers, func(context *Context) {
				// Check if file exists and/or if we have permission to access it
				if _, err := fs.Open(ctx.Path[len(first)+1:]); err != nil {
					context.Fail("not found")
					return
				}
				fmt.Println("found file")
				fileServer.ServeHTTP(context.Res, context.Req)
			})
		}
	}
	if !flag {
		node, params := router.parseTrieRoute(ctx.Method, searchParts)
		if node != nil {
			ctx.Params = params
			key := ctx.Method + "-" + node.fullPath
			handler := router.handlerFuncMap[key]
			ctx.handlers = append(ctx.handlers, handler) // 路由匹配的handler放在最后
		} else {
			ctx.handlers = append(ctx.handlers, func(context *Context) {
				context.Send("Not Found", 404)
			})
		}
	}
	ctx.Next()
}

func (router *Router) parseTrieRoute(method string, searchParts []string) (*TrieNode, map[string]string) {
	root, ok := router.rootNode[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string)
	node := root.searchPath(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.fullPath)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
		}
		return node, params
	}
	return nil, nil
}

func (router *Router) printRouter() {
	// bfs 打印当前trie树结构
	for key, val := range router.rootNode {
		// method root
		fmt.Println(key)
		// bfs
		type Item struct {
			node   *TrieNode
			height int
		}
		queue := make([]Item, 0)
		queue = append(queue, Item{
			node:   val,
			height: 0,
		})
		preHeight := 0
		for {
			if len(queue) == 0 {
				break
			}
			cur := queue[0]   // 入队
			queue = queue[1:] // 出队
			if cur.height != preHeight {
				fmt.Printf("\n")
				preHeight++
			}
			fmt.Printf("%s ", cur.node.curPath)
			for _, item := range cur.node.childNodes {
				queue = append(queue, Item{
					node:   item,
					height: cur.height + 1,
				})
			}
		}
		fmt.Printf("\n-----------------\n")
	}
}
