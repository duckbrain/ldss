package main

import (
	"ldss/lib"
	"log"
	"os"
)

type app interface {
	run()
	setInfo(args []string)
}

type appinfo struct {
	args  []string
	fmt   *log.Logger
	efmt  *log.Logger
	debug *log.Logger
	lang  *lib.Language
}

func (a *appinfo) setInfo(args []string) {
	a.args = args
	a.fmt = log.New(os.Stdin, "", 0)
	a.efmt = log.New(os.Stderr, "", 0)
	a.debug = a.getDebug()
	if lang, err := lib.LookupLanguage(Config().Get("Language").(string)); err != nil {
		panic(err)
	} else {
		a.lang = lang
	}
}

var apps map[string]app

func addApp(name string, a app) {
	if apps == nil {
		apps = make(map[string]app)
	}
	apps[name] = a
}
