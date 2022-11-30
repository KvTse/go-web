package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	/**
	所有路由和处理器的映射
	*/
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRouter(method, pattern, handler)
}
func (engine *Engine) RegisterGetRouter(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) RegisterPostRouter(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}
func (engine *Engine) StartServer(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
