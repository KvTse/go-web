package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

/*
*
程序启动入口
*/
func main() {

	engine := gee.New()
	engine.Use(gee.Logger())

	engine.RegisterGetRouter("/hello/:name", func(c *gee.Context) {
		c.WriteString(http.StatusOK, "GET hello...%s", c.Params["name"])
	})
	engine.RegisterGetRouter("/hello/star/*", func(c *gee.Context) {
		c.WriteString(http.StatusOK, "GET hello...%s", c.Params["name"])
	})

	group := engine.Group("/v1")
	{
		group.RegisterGetRouter("/hello", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>hello v1 group</h1>")
		})
		group.RegisterGetRouter("/", func(c *gee.Context) {
			c.WriteString(http.StatusOK, "I'm the v1 root")
		})
	}

	group2 := engine.Group("/v2")
	{
		group2.RegisterGetRouter("/testMiddlewares", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1> testMiddlewares </h1>")
		})
	}
	group2.Use(onlyForV2RGroup())
	err := engine.StartServer(":8099")
	if err != nil {
		panic("fail to start server...")
	}
}

func onlyForV2RGroup() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error...")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
