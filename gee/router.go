package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	/**
	路由路径和处理器的映射
	*/
	handlers map[string]HandlerFunc
}

// NewRouter /*新建router
func newRouter() *router {
	return &router{roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

/*
*
添加路由
POST-/login -> handler
*/
func (router *router) addRouter(methodType string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	_, ok := router.roots[methodType]
	if !ok {
		// 创建节点
		router.roots[methodType] = &node{}
	}
	router.roots[methodType].insert(pattern, parts, 0)

	key := methodType + "-" + pattern
	log.Printf("Router %s", key)
	router.handlers[key] = handler
}
func (router *router) getRouter(methodType string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := router.roots[methodType]
	if !ok {
		return nil, nil
	}
	searchedNode := root.search(searchParts, 0)
	if searchedNode != nil {
		parts := parsePattern(searchedNode.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return searchedNode, params
	}
	return nil, nil
}
func (router *router) getRouters(methodType string) []*node {
	root, ok := router.roots[methodType]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}
func (router *router) handle(c *Context) {
	node, params := router.getRouter(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		router.handlers[key](c)
	} else {
		c.WriteString(http.StatusNotFound, "404 NOT FOUND :%s\n", c.Path)
	}
}
