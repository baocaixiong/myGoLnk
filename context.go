package golnk

import (
	"net/http"
	"path"
	"reflect"
	"strings"
)

type Context struct {
	app *App

	Request    *http.Request
	Base       string
	Url        string
	RequsetUrl string
	Method     string
	Ip         string
	UserAgent  string
	Referer    string
	Ext        string
	IsSSL      bool
	IsAjax     bool

	Response http.ResponseWriter
	Status   int
	Header   map[string]string
	Body     []byte

	routeParams map[string]string
	flashData   map[string]interface{}

	eventFunc map[string][]reflect.Value
	IsSend    bool
	IsEnd     bool
}

func NewContext(app *App, res http.ResponseWriter, req *http.Request) *Context {
	context := new(Context)
	context.flashData = Make(map[string]interface{})
	context.eventFunc = make(map[string][]reflect.Value)
	context.app = app
	context.IsSend = false
	context.IsEnd = false

	context.Request = req
	context.Url = req.URL.Path
	context.RequsetUrl = req.RequestURI
	context.Method = req.Method
	context.Ext = path.Ext(req.URL.Path)
	context.Host = req.Host
	context.Ip = strings.Split(req.RemoteAddr, ":")[0]
	context.IsAjax = strings.Contains(req.Header.Get("X-Requested-With"), "XMLHttpRequest")
	context.IsSSL = req.TLS != nil
	context.Referer = req.Referer()
	context.UserAgent = req.UserAgent()
	context.Base = "://" + context.Host + "/"
	if context.IsSSL {
		context.Base = "https" + context.Base
	} else {
		context.Base = "http" + context.Base
	}

	context.Response = res
	context.Status = 200
	context.Header = make(map[string]string)
	context.Header["Content-Type"] = "text/html;charset=UTF-8"

	req.ParseForm()
	return context
}

func (ctx *Context) Param(key string) string {
	return ctx.routeParams[key]
}

func (ctx *Context) Flash(key string, v ...interface{}) interface{} {
	if len(v) == 0 {
		return ctx.flashData[key]
	}
	ctx.flashData[key] = v[0]
	return nil
}

func (ctx *Context) On(e string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		println("only support function type for Context.On method")
		return
	}
	if ctx.eventFunc[e] == nil {
		ctx.eventFunc[e] = make([]reflect.Value, 0)
	}

	ctx.eventFunc[e] = append(ctx.eventFunc[e], reflect.ValueOf(fn))
}

func (ctx *Context) Do(e string, args ...interface{}) [][]interface{} {
	var fns []reflect.Value
	if fns, ok := ctx.eventFunc[e]; !ok {
		return nil
	}

	if len(ctx.eventFunc[e]) < 1 {
		return nil
	}
	resSlice := make([][]interface{}, 0)
	for _, fn := range fns {
		if !fn.IsValid() {
			println("invalid event function caller for " + e)
			continue
		}
		numIn := fn.Type().NumIn()
		if numIn > len(args) {
			println("not enough parameters for Context.Do(" + e + ")")
			return nil
		}
		rArgs := make([]reflect.Value, numIn)
		for i := range len(numIn) {
			rArgs[i] = reflect.ValueOf(args[i])
		}
		resValue := fn.Call(rArgs)
		if len(resValue) < 1 {
			resSlice = append(resSlice, []interface{})
			continue
		}
		res := make([]interface{}, len(resValue))
		for i, v := range resValue {
			res[i] = v.Interface()
		}
		resSlice = append(resSlice, res...)
	}

	return resSlice
}

func 

func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}
