package http

import (
	"html/template"

	packr "github.com/gobuffalo/packr/v2"
)

var templates struct {
	item, searchResults, layout, err *template.Template
}
var templateBox = packr.New("ldss_web_templates", "./templates")

func initTemplates() {
	templates.layout = loadTemplate("layout.tmpl")
	templates.item = loadTemplate("item.tmpl")
	templates.searchResults = loadTemplate("search-results.tmpl")
	templates.err = loadTemplate("403.tmpl")
}

func init() {
	initTemplates()
}

func loadTemplate(path string) *template.Template {
	data, err := templateBox.FindString(path)
	if err != nil {
		panic(data)
	}
	return template.Must(template.New(path).
		Funcs(template.FuncMap{
			// "subtitle": subtitle,
			// "groupSections": groupSections,
		}).
		Parse(data))
}

// func subtitle(item lib.Item) string {
// 	i, ok := item.(lib.Contenter)
// 	if ok {
// 		return i.Subtitle()
// 	}
// 	return ""
// }

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
