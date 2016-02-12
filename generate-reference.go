// Generates reference files for a given language

// +build !release

package main

import (
	"ldss/lib"
)

type generateReference struct {
	appinfo
}

func init() {
	addApp("generate-reference", &generateReference{})
}

func (app *generateReference) run() {
	langId := app.args[1]
	app.efmt.Println(langId)
	lang, err := lib.LookupLanguage(langId)
	if err != nil {
		panic(err)
	}
	app.efmt.Println(lang.String())
	messages := lib.DownloadAll(lang, false)
	for m := range messages {
		app.efmt.Println(m.String())
	}
	catalog, err := lang.Catalog()
	if err != nil {
		panic(err)
	}

	for _, b := range catalog.Books() {
		app.runBook(b)
	}
	for _, f := range catalog.Folders() {
		app.runFolder(f)
	}
}

func (app *generateReference) runFolder(f *lib.Folder) {
	app.fmt.Printf("%v:%v:%v\n", f.ID(), f.Name(), f.Path())
	for _, b := range f.Books() {
		app.runBook(b)
	}
	for _, f := range f.Folders() {
		app.runFolder(f)
	}
}

func (app *generateReference) runBook(b *lib.Book) {
	app.fmt.Printf("%v:%v\n", b.Name(), b.Path())
	nodes, err := b.Index()
	if err != nil {
		app.efmt.Println(err)
		return
	}
	for _, n := range nodes {
		app.runNode(n)
	}
}

func (app *generateReference) runNode(n *lib.Node) {

}
