package golnk

import (
	goUrl "net/url"
	"path"
	"regexp"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
	PATCH  = "PATCH"
	OPTION = "OPTIONE"
	HEAD   = "HEAD"
)

type Router struct {
	routeSlice []*Route
	Get        func(pattern string, fn ...Handler)
	Post       func(pattern string, fn ...Handler)
	Delete     func(pattern string, fn ...Handler)
	Put        func(pattern string, fn ...Handler)
	Patch      func(pattern string, fn ...Handler)
	Option     func(pattern string, fn ...Handler)
	Head       func(pattern string, fn ...Handler)
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
	rt := new(Router)
	rt.routeSlice = make([]*Route, 0)
	rt.Get = rt.makeFunc(GET)
	rt.Post = rt.makeFunc(POST)
	rt.Delete = rt.makeFunc(DELETE)
	rt.Put = rt.makeFunc(PUT)
	rt.Patch = rt.makeFunc(PATCH)
	rt.Option = rt.makeFunc(OPTION)
	rt.Head = rt.makeFunc(HEAD)
	return rt
}

func newRoute() *Route {
	route := new(Route)
	route.params = make([]string, 0)
	return route
}

func (rt *Router) makeFunc(method string) func(pattern string, fn ...Handler) {
	return func(pattern string, fn ...Handler) {
		route := newRoute()
		route.regex, route.params = rt.parsePattern(pattern)
		route.method = method
		route.fn = fn
		rt.routeSlice = append(rt.routeSlice, route)
	}
}

func (rt *Router) AddFunc(method string, pattern string, fn ...Handler) {
	methods := strings.Split(method, ",")
	for _, m := range methods {
		switch strings.Trim(m, " ") {
		case "GET":
			rt.Get(pattern, fn...)
		case "POST":
			rt.Post(pattern, fn...)
		case "PUT":
			rt.Put(pattern, fn...)
		case "DELETE":
			rt.Delete(pattern, fn...)
		case "PATCH":
			rt.Patch(pattern, fn...)
		case "HEAD":
			rt.Head(pattern, fn...)
		case "OPTION":
			rt.Option(pattern, fn...)
		default:
			println("unknow route method: " + method)
		}
	}
}

func (rt *Router) parsePattern(pattern string) (regex *regexp.Regexp, params []string) {
	params = make([]string, 0)
	segments := strings.Split(goUrl.QueryEscape(pattern), "%2F")
	for i, v := range segments {
		if strings.HasPrefix(v, "%3A") {
			segments[i] = `([\w-%]+)`
			params = append(params, strings.TrimPrefix(v, "%3A"))
		}
	}
	regex, _ = regexp.Compile("^" + strings.Join(segments, "/") + "$")
	return
}

func (rt *Router) Find(url string, method string) (params map[string]string, fn []Handler) {
	sfx := path.Ext(url)
	url = strings.Replace(url, sfx, "", -1)
	url = goUrl.QueryEscape(url) // : => %3A
	if !strings.HasSuffix(url, "%2F") && sfx == "" {
		url += "%2F" //加上 /
	}

	url = strings.Replace(url, "%2F", "/", -1)
	for _, r := range rt.routeSlice {
		p := r.regex.FindStringSubmatch(url)
		if len(p) != len(r.params)+1 {
			continue
		}
		params = make(map[string]string)
		for i, n := range r.params {
			params[n] = p[i+1]
		}
		fn = r.fn
		return
	}
	return nil, nil
}
