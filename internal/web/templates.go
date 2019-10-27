package web

import (
	"html/template"

	"github.com/duckbrain/ldss/lib"
	packr "github.com/gobuffalo/packr/v2"
)

var templates *webtemplates
var templateBox = packr.New("ldss_web_templates", "./templates")

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
	data, err := templateBox.FindString(path)
	if err != nil {
		panic(data)
	}
	return template.Must(template.New(path).
		Funcs(template.FuncMap{
			"subtitle": subtitle,
			// "groupSections": groupSections,
		}).
		Parse(data))
}

func subtitle(item lib.Item) string {
	i, ok := item.(lib.Contenter)
	if ok {
		return i.Subtitle()
	}
	return ""
}

// type groupedSections map[string][]lib.Contenter

// func groupSections(items []lib.Item) groupedSections {
// 	nodeMap := make(groupedSections)
// 	for _, item := range items {
// 		contenter, ok := item.(lib.Contenter)
// 		if !ok {
// 			return nil
// 		}
// 		key := contenter.SectionName()
// 		list, ok := nodeMap[key]
// 		if ok {
// 			list = append(list, contenter)
// 		} else {
// 			list = []lib.Contenter{contenter}
// 		}
// 		nodeMap[key] = list
// 	}
// 	return nodeMap
// }
