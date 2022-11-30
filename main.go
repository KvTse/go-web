package main

import (
	"gee"
	"net/http"
)

/*
*
程序启动入口
*/
func main() {
	engine := gee.New()
	engine.RegisterGetRouter("/", func(c *gee.Context) {
		c.WriteString(http.StatusOK, "success")
	})
	engine.RegisterPostRouter("/hello", func(c *gee.Context) {
		name := c.FormValue("name")
		c.WriteString(http.StatusOK, "hello...%s", name)
	})
	err := engine.StartServer(":8099")
	if err != nil {
		panic("fail to start server...")
	}
}
