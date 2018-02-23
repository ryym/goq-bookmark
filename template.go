package main

import (
	"bytes"
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func NewTemplate(globPath string) *Template {
	return &Template{
		templates: template.Must(template.ParseGlob(globPath)),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, "main", &TemplateContent{name, data, t.templates})
}

type TemplateContent struct {
	Name      string
	Data      interface{}
	templates *template.Template
}

func (tc *TemplateContent) Content(data interface{}) template.HTML {
	buf := new(bytes.Buffer)
	err := tc.templates.ExecuteTemplate(buf, tc.Name, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(string(buf.Bytes()))
}
