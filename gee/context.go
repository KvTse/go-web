package gee

import (
	"fmt"
	"net/http"
)

// Context 封装一次http请求相关信息
type Context struct {
	Write      http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newContext(write http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Write:  write,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
func (c *Context) SetHeader(key string, value string) {
	c.Write.Header().Set(key, value)
}
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Write.WriteHeader(code)
}
func (c *Context) WriteString(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Write.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		return
	}
}

func (c *Context) FormValue(key string) string {
	return c.Req.FormValue(key)
}
