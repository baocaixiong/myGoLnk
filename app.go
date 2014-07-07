package golnk

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type App struct {
	router  *Router
	routerC map[string]*routerCache
	view    *View
	middle  []Handler
	inter   map[string]Handler
	config  *Config
	Get     func(key string, fn ...Handler)
	Post    func(key string, fn ...Handler)
	Delete  func(key string, fn ...Handler)
	Put     func(key string, fn ...Handler)
	Patch   func(key string, fn ...Handler)
	Option  func(key string, fn ...Handler)
	Head    func(key string, fn ...Handler)
}

func NewApp() *App {
	a := new(App)
	a.router = NewRouter()
	a.routerC = make(map[string]*routerCache)
	a.middle = make([]Handler, 0)
	a.inter = make(map[string]Handler)
	a.config, _ = NewConfig("config.json")
	a.view = NewView(a.config.StringOr("view_dir", "view"))

	a.Get = a.makeRouterFunc(GET)
	a.Post = a.makeRouterFunc(POST)
	a.Delete = a.makeRouterFunc(DELETE)
	a.Put = a.makeRouterFunc(PUT)
	a.Patch = a.makeRouterFunc(PATCH)
	a.Option = a.makeRouterFunc(OPTION)
	a.Head = a.makeRouterFunc(HEAD)

	return a
}

// 添加处理中间件
func (app *App) Use(h ...Handler) {
	app.middle = append(app.middle, h...)
}

func (app *App) Config() *Config {
	return app.config
}

func (app *App) View() *View {
	return app.view
}

func (app *App) handler(res http.ResponseWriter, req *http.Request) {
	context := NewContext(app, res, req)
	defer func() {
		e := recover()
		if e == nil {
			context = nil
			return
		}
		context.Body = []byte(fmt.Sprint(e))
		context.Status = 503
		println(string(context.Body))
		debug.PrintStack()
		if _, ok := app.inter["recover"]; ok {
			app.inter["recover"](context)
		}
		if !context.IsEnd {
			context.End()
		}
		context = nil
	}()

	if _, ok := app.inter["static"]; ok {
		app.inter["static"](context)
		if context.IsEnd {
			return
		}
	}

	if len(app.middle) > 0 {
		for _, h := range app.middle {
			h(context)
			if context.IsEnd {
				break
			}
		}
	}

	if context.IsEnd {
		return
	}

	var (
		params map[string]string
		fn     []Handler
		url    = req.URL.Path
	)

	if _, ok := app.routerC[url]; ok {
		params = app.routerC[url].param
		fn = app.routerC[url].fn
	} else {
		params, fn = app.router.Find(url, req.Method)
	}

	if params != nil && fn != nil {
		context.routerParams = params
		rc := new(routerCache)
		rc.param = params
		rc.fn = fn
		app.routerC[url] = rc
		for _, f := range fn {
			f(context)
			if context.IsEnd {
				break
			}
		}

		if !context.IsEnd {
			context.End()
		}
	} else {
		println("router is missing at " + req.URL.Path)
		context.Status = 404
		if _, ok := app.inter["notFound"]; ok {
			app.inter["notFound"](context)
			if !context.IsEnd {
				context.End()
			}
		} else {
			context.Throw(404)
		}
	}
	context = nil
}

func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.handler(res, req)
}

func (app *App) Run() {
	addr := app.config.StringOr("app.server", "localhost:9001")
	println("http server run at: " + addr)
	e := http.ListenAndServe(addr, app)
	panic(e)
}

func (app *App) Route(method string, key string, fn ...Handler) {
	app.router.AddFunc(method, key, fn...)
}

func (app *App) makeRouterFunc(method string) func(key string, fn ...Handler) {
	return func(key string, fn ...Handler) {
		app.router.AddFunc(method, key, fn...)
	}
}

func (app *App) AddInter(key string, h Handler) {
	app.inter[key] = h
}
