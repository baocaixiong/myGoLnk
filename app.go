package golnk

import ()

type App struct {
	router  *Router
	routerC map[string]*routerCache
	view    *View
	middle  []Handler
	inter   map[string]Handler
	config  *Config
}

func NewApp() *App {
	a := new(App)
	a.router = NewRouter()
	a.routerC = make(map[string]*routerCache)
	a.middle = make([]Handler, 0)
	a.inter = make(map[string]Handler)
	a.config, _ = NewConfig("config.json")
	a.view = NewView(a.config.String("view_dir", "view"))
	println(a.config.String("nihao.haha", "123"))

	return a
}
