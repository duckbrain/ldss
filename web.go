package main

import (
	"strings"
	"fmt"
	"net/http"
	"log"
	"os"
	"ldslib"
)

type web struct {
	args []string
	config Config
}

func (app web) run() {
	efmt := log.New(os.Stderr, "", 0)
	efmt.Printf("Starting web server on port: %v\n", app.config.op.WebPort)
	
	http.HandleFunc("/", app.handler)
	http.ListenAndServe(fmt.Sprintf(":%v", app.config.op.WebPort), nil)
}

func (app *web) handler(w http.ResponseWriter, r *http.Request) {
	lang, err := app.config.Library.Language(r.URL.Query().Get("lang"))
	if err != nil {
		lang = app.config.SelectedLanguage()
	}
	catalog := app.config.Catalog(lang)
	
	switch strings.ToLower(r.URL.Path) {
		case "/download":
		switch r.Method {
			case "GET":
				fmt.Fprintf(w, "Download Page Here")
			case "POST":
				fmt.Fprintf(w, "Download content here")
		}
		case "/":
			app.print(w, r, catalog)
		case "/favicon.ico":
			w.Header().Set("Content-type", "image/x-icon")
			data, err := Asset("data/favicon.ico")
			if err != nil {
				panic (err)
			}
			w.Write(data)
		default:
			item, err := app.config.Library.Lookup(r.URL.Path, catalog)
			if err != nil {
				panic(err)
			}
			app.print(w, r, item)
	}
	
}

func (app *web) print(w http.ResponseWriter, r *http.Request, item ldslib.CatalogItem) {
	
	fmt.Fprintf(w, `
		<html>
		<head>
			<link rel="icon" href="/favicon.ico" sizes="16x16 32x32 128x128" type="image/vnd.microsoft.icon">
			<title>%v</title>
		</head>
		<body>`, item)
	fmt.Fprintf(w, "<h1>%v</h1>", item)
	
	linkQuery := fmt.Sprintf("?lang=%v", item.Language().GlCode)
	
	switch item.(type) {
		case ldslib.Node:
			node := item.(ldslib.Node)
			if node.HasContent {
				content, _ := app.config.Library.Content(node)
				fmt.Fprintf(w, content)
			} else {
				children, err := app.config.Library.Children(item)
				if err != nil {
					panic(err)
				}
				fmt.Fprint(w, "<ul>")
				for _, child := range children {
					fmt.Fprintf(w, `<li><a href="%v%v">%v</a></li>`, child.Path(), linkQuery, child.DisplayName())
				}
				fmt.Fprint(w, "</ul>")
			}
		default:
			children, err := app.config.Library.Children(item)
			if err != nil {
				panic(err)
			}
			fmt.Fprint(w, "<ul>")
			for _, child := range children {
				fmt.Fprintf(w, `<li><a href="%v%v">%v</a></li>`, child.Path(), linkQuery, child.DisplayName())
			}
			fmt.Fprint(w, "</ul>")
	}
	fmt.Fprint(w, "</body>")
}