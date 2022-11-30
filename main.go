package main

import (
	"fmt"
	"gee"
	"net/http"
)

/*
*
程序启动入口
*/
func main() {
	engine := gee.New()
	engine.RegisterGetRouter("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "request success\ndo something of get method...")
	})
	engine.RegisterPostRouter("/hello", func(writer http.ResponseWriter, request *http.Request) {
		name := request.FormValue("name")
		fmt.Fprintf(writer, "hello...%s", name)
	})
	err := engine.StartServer(":8099")
	if err != nil {
		panic("fail to start server...")
	}
}
