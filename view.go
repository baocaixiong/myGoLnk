package golnk

import (
	"html/template"
)

type View struct {
	// template directory
	Dir string
	// view functions map
	FuncMap template.FuncMap
	// Cache Flag
	IsCache bool
	// template cache map
	templateCache map[string]*template.Template
}

func NewView(dir string) *View {
	v := new(View)
	v.Dir = dir
	v.FuncMap = make(template.FuncMap)
	v.FuncMap["Html"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	v.IsCache = false
	v.templateCache = make(map[string]*template.Template)
	return v
}
