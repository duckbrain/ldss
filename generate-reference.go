// Generates reference files for a given language

// +build !release

package main

import "ldss/lib"

type generateReference struct {
	appinfo
}

func init() {
	addApp("generate-reference", &generateReference{})
}

func (app *generateReference) run() {
	lang, err := lib.LookupLanguage(app.args[0])
	if err != nil {
		panic(err)
	}
	app.fmt.Println(lang.String())
}
