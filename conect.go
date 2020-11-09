package main

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"path/filepath"
)

type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request *http.Request
}

type HandlerFunc func(*Context)


// templates: 템플릿 객체를 보관하기 위한 map
var templates = map[string]*template.Template{}


func (c *Context) RenderJson(v interface{}){
	// HTTP Status를 StatusOK로 지정
	c.ResponseWriter.WriteHeader(http.StatusOK)
	// Content-Type을 Allplication/json
	c.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		// 에러 발생 시 RenderErr 메세드 호출
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderXml(v interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}

}


func (c *Context) RenderErr(code int, err error) {
	if err != nil {
		if code > 0 {
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			defaltErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaltErr), defaltErr)
		}
	}
}

func (c *Context) RenderTemplate(path string, v interface{}) {
	t, ok := templates[path]
	if !ok {
		t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		templates[path] = t

	}
	t.Execute(c.ResponseWriter, v)
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}