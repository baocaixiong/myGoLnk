package golnk

import (
	"net/http"
)

type Context struct {
	app *App
	req *http.Request
	res http.ResponseWriter
}

func NewContext(app *App, res http.ResponseWriter, req *http.Request) *Context {
	context := new(Context)

}
