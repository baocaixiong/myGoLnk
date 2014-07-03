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
	a.view = NewView(a.config.StringOr("view_dir", "view"))

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
