package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	/**
	所有路由和处理器的映射
	*/
	router *router
	*RouterGroup
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}
func (engine *Engine) LoadHtmlGlob(pattern string) {
	engine.htmlTemplates = template.Must(
		template.New("").
			Funcs(engine.funcMap).
			ParseGlob(pattern))
}

type RouterGroup struct {
	prefix      string // 前缀
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 创建一个分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
		parent: group}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// {(engine *Engine) -> (group *RouterGroup)} addRouter一下子就变成是RouterGroup的方法了
func (group *RouterGroup) addRouter(method string, prefix string, handler HandlerFunc) {
	pattern := group.prefix + prefix
	group.engine.router.addRouter(method, pattern, handler)
}
func (group *RouterGroup) RegisterGetRouter(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) RegisterPostRouter(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

func (engine *Engine) StartServer(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutionPath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutionPath, http.FileServer(fs))
	return func(c *Context) {
		filePath := c.Param("filePath")
		if _, err := fs.Open(filePath); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
func (group *RouterGroup) Static(relativePath string, root string) {
	staticHandler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filePath")
	group.RegisterGetRouter(urlPattern, staticHandler)
}
