package web

import (
	"html/template"

	"github.com/duckbrain/ldss/assets"
	"github.com/duckbrain/ldss/lib"
)

var templates *webtemplates

type webtemplates struct {
	nodeChildren, nodeContent, searchResults, layout, err *template.Template
}

func initTemplates() {
	templates = &webtemplates{}
	templates.layout = loadTemplate("layout.tpl")
	templates.nodeContent = loadTemplate("node-content.tpl")
	templates.nodeChildren = loadTemplate("node-children.tpl")
	templates.searchResults = loadTemplate("search-results.tpl")
	templates.err = loadTemplate("403.tpl")
}

func loadTemplate(path string) *template.Template {
	data, err := assets.Asset("data/web/templates/" + path)
	if err != nil {
		panic(err)
	}
	temp := template.New(path)
	temp.Funcs(template.FuncMap{
		"subtitle": func(item lib.Item) string {
			node, ok := item.(*lib.Node)
			if ok {
				return node.Subtitle
			} else {
				return ""
			}
		},
	})
	temp, err = temp.Parse(string(data))
	if err != nil {
		panic(err)
	}
	return temp
}
