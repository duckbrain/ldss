// Generates reference files for a given language

// +build debug
// +build exclude

package main

import (
	"ldss/lib"
	"strconv"
	"strings"
)

type generateReference struct {
	appinfo
	lang *lib.Language
	cat  *lib.Catalog
}

func init() {
	addApp("generate-reference", &generateReference{})
}

func (app generateReference) register(*Configuration) {}

func (app *generateReference) lookup(path string) lib.Item {
	item, err := app.cat.LookupPath(path)
	if err != nil {
		panic(err)
	}
	return item
}

func (app *generateReference) run() {
	if len(app.args) != 2 {
		panic("Invalid number of arguments, expects language id")
	}
	langId := app.args[1]
	app.efmt.Println(langId)
	var err error
	if app.lang, err = lib.LookupLanguage(langId); err != nil {
		panic(err)
	}
	if app.cat, err = app.lang.Catalog(); err != nil {
		panic(err)
	}
	app.efmt.Println(app.lang.String())
	messages := lib.DownloadAll(app.lang, false)
	for m := range messages {
		app.efmt.Println(m.String())
	}

	app.runScriptureVolume("/scriptures/ot")
	app.runScriptureVolume("/scriptures/nt")
	app.runScriptureVolume("/scriptures/bofm")
	app.runScriptureVolume("/scriptures/pgp")
	app.runDandC(app.lookup("/scriptures/dc-testament").(*lib.Book))
}

// Generates lookup names from user readable strings
func (app *generateReference) userNames(name string) []string {
	name = strings.ToLower(name)
	return []string{name}
}

// Generates lookup names from the last component of a path
func (app *generateReference) pathNames(name string) []string {
	name = name[strings.LastIndex(name, "/")+1:]
	name = strings.ToLower(name)
	name = strings.Replace(name, "-", " ", 100)
	return []string{name}
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func (app *generateReference) genSimple(n lib.Item, hash string) {
	app.gen(append(app.userNames(n.Name()), app.pathNames(n.Path())...), n.Path()+hash)
}

func (app *generateReference) genParent(n lib.Item, hash string) {
	p := n.Parent()
	names := []string{}
	pnames := append(app.userNames(p.Name()), app.pathNames(p.Path())...)
	nnames := append(app.userNames(n.Name()), app.pathNames(n.Path())...)
	for _, pname := range pnames {
		for _, nname := range nnames {
			names = append(names, pname+" "+nname)
		}
	}
	app.gen(names, n.Path())
}

func (app *generateReference) gen(matches []string, path string) {
	cleaned := []string{}
	for _, value := range matches {
		if !stringInSlice(value, cleaned) {
			cleaned = append(cleaned, value)
		}
	}
	app.fmt.Printf("%v:%v\n", strings.Join(cleaned, ":"), path)
}

func (app *generateReference) comment(comment string) {
	app.fmt.Printf("#%v\n", comment)
}

func (app *generateReference) runScriptureVolume(path string) {
	b, err := app.cat.LookupPath(path)
	if err != nil {
		return
	}
	app.fmt.Println("")
	app.comment(b.Name())

	app.genSimple(b, "")
	nodes, err := b.Children()
	if err != nil {
		app.efmt.Println(err)
		return
	}
	for _, n := range nodes {
		app.runScriptureBook(n)
	}
}

func (app *generateReference) runScriptureBook(n lib.Item) {
	if _, err := app.cat.LookupPath(n.Path() + "/1"); err == nil {
		if _, err := app.cat.LookupPath(n.Path() + "/2"); err != nil {
			// Is a single chapter book
			names := app.userNames(n.Name())
			names = append(names, app.userNames(n.Name()+" 1")...)
			names = append(names, app.pathNames(n.Path())...)
			names = append(names, app.pathNames(n.Path()+"-1")...)
			app.gen(names, n.Path())
		} else {
			// Is a multiple chapter book
			app.genSimple(n, "#")
		}
	} else {
		app.genParent(n, "")
	}
}

func (app *generateReference) runDandC(n lib.Item) {
	app.fmt.Println("")
	app.comment(n.Name())
	app.genSimple(n, "#")
	//TODO Generate the number regex
	children, err := n.Children()
	if err != nil {
		panic(err)
	}
	for _, c := range children {
		path := c.Path()
		path = path[strings.LastIndex(path, "/")+1:]
		if _, err := strconv.Atoi(path); err != nil {
			// Print the ones that don't end in numbers
			app.genParent(c, "")
		}
	}
}
