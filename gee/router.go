package gee

import (
	"log"
	"net/http"
)

type router struct {
	/**
	路由路径和处理器的映射
	*/
	handlers map[string]HandlerFunc
}

// NewRouter /*新建router
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

/*
*
添加路由
POST-/login -> handler
*/
func (router *router) addRouter(methodType string, pattern string, handler HandlerFunc) {
	key := methodType + "-" + pattern
	log.Printf("Router %s", key)
	router.handlers[key] = handler
}
func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c)
	} else {
		c.WriteString(http.StatusNotFound, "404 NOT FOUND :%s\n", c.Path)
	}
}
