package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

/*
*
程序启动入口
*/
//func main() {
//
//	engine := gee.New()
//	engine.Use(gee.Logger())
//
//	engine.RegisterGetRouter("/hello/:name", func(c *gee.Context) {
//		c.WriteString(http.StatusOK, "GET hello...%s", c.Params["name"])
//	})
//	engine.RegisterGetRouter("/hello/star/*", func(c *gee.Context) {
//		c.WriteString(http.StatusOK, "GET hello...%s", c.Params["name"])
//	})
//
//	group := engine.Group("/v1")
//	{
//		group.RegisterGetRouter("/hello", func(c *gee.Context) {
//			c.HTML(http.StatusOK, "<h1>hello v1 group</h1>")
//		})
//		group.RegisterGetRouter("/", func(c *gee.Context) {
//			c.WriteString(http.StatusOK, "I'm the v1 root")
//		})
//	}
//
//	group2 := engine.Group("/v2")
//	{
//		group2.RegisterGetRouter("/testMiddlewares", func(c *gee.Context) {
//			c.HTML(http.StatusOK, "<h1> testMiddlewares </h1>")
//		})
//	}
//	group2.Use(onlyForV2RGroup())
//
//	engine.Static("/assets", "./static")
//
//	err := engine.StartServer(":8099")
//	if err != nil {
//		panic("fail to start server...")
//	}
//}

func main() {

	engine := gee.New()
	engine.Use(gee.Logger())

	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHtmlGlob("templates/*")
	engine.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	engine.RegisterGetRouter("/", func(c *gee.Context) {
		c.Html(http.StatusOK, "css.tmpl", gee.H{})
	})
	engine.RegisterGetRouter("/students", func(c *gee.Context) {
		c.Html(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	engine.RegisterGetRouter("/date", func(c *gee.Context) {
		c.Html(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	err := engine.StartServer(":8099")
	if err != nil {
		panic("fail to start server...")
	}
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV2RGroup() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error...")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
