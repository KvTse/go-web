package gee

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	/**
	所有路由和处理器的映射
	*/
	router map[string]HandlerFunc
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handlerFunc, ok := engine.router[key]; ok {
		handlerFunc(w, req)
	} else {
		_, _ = fmt.Fprintf(w, "404 url not found -> %s", req.URL)
	}

}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

/*
*
添加路由
POST-/login -> handler
*/
func (engine *Engine) addRouter(methodType string, pattern string, handler HandlerFunc) {
	key := methodType + "-" + pattern
	log.Printf("Router %s", key)
	engine.router[key] = handler
}

// RegisterGetRouter /** 暴露添加GET路由的方法
func (engine *Engine) RegisterGetRouter(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}
func (engine *Engine) RegisterPostRouter(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}
func (engine *Engine) StartServer(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
