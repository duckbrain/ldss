package main

import (
	"html/template"
)

type webtemplates struct {
	nodeChildren, nodeContent, layout, err *template.Template
}

type webtemplateData struct {
	data interface{}
}

func (app *web) initTemplates() {
	app.templates = &webtemplates{}
	app.templates.layout = app.loadTemplate("layout.tpl")
	app.templates.nodeContent = app.loadTemplate("node-content.tpl")
	app.templates.nodeChildren = app.loadTemplate("node-children.tpl")
	app.templates.err = app.loadTemplate("403.tpl")
}

func (app *web) loadTemplate(path string) *template.Template {
	data, err := Asset("data/web/templates/" + path)
	if err != nil {
		panic(err)
	}
	temp := template.New(path)
	temp, err = temp.Parse(string(data))
	if err != nil {
		panic(err)
	}
	return temp
}

func (app *web) loadLayoutTemplate(path string, layout *template.Template) {
	//temp := app.loadTemplate(path)
	// Need to create Executor interface to make this work
}
