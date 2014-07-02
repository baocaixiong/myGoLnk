package golnk

import (
	"regexp"
)

type Router struct {
	routeSlice []*Route
}

type Route struct {
	regex  *regexp.Regexp
	method string
	params []string
	fn     []Handler
}

type Handler func(context *Context)

type routerCache struct {
	param map[string]string
	fn    []Handler
}

func NewRouter() *Router {
	router := new(Router)
	router.routeSlice = make([]*Route, 0)
	return router
}
