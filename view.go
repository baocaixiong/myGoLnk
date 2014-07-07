package golnk

import (
	"bytes"
	"html/template"
	"path"
	"strings"
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

func (v *View) getTemplateInstance(tpl []string) (*template.Template, error) {
	key := strings.Join(tpl, "-")
	if v.IsCache {
		if v.templateCache[key] != nil {
			return v.templateCache[key], nil
		}
	}
	var (
		t    *template.Template
		e    error
		file []string = make([]string, len(tpl))
	)
	for i, tp := range tpl {
		file[i] = path.Join(v.Dir, tp)
	}
	t = template.New(path.Base(tpl[0]))
	t.Funcs(v.FuncMap)
	t, e = t.ParseFiles(file...)
	if e != nil {
		return nil, e
	}
	if v.IsCache {
		v.templateCache[key] = t
	}
	return t, nil
}

func (view *View) Render(tpl string, data map[string]interface{}) ([]byte, error) {
	t, e := view.getTemplateInstance(strings.Split(tpl, ","))
	if e != nil {
		return nil, e
	}

	var buf bytes.Buffer
	e = t.Execute(&buf, data)
	if e != nil {
		return nil, e
	}

	return buf.Bytes(), nil
}

func (v *View) NOCache() {
	v.IsCache = false
	v.templateCache = make(map[string]*template.Template)
}
