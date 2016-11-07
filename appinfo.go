package main

import (
	"ldss/lib"
	"log"
	"os"
)

type app interface {
	run()
	register(config *Configuration)
	setInfo(config *Configuration, args []string)
}

type appinfo struct {
	args   []string
	config *Configuration
	fmt    *log.Logger
	efmt   *log.Logger
	debug  *log.Logger
	lang   *lib.Language
}

func (a *appinfo) setInfo(config *Configuration, args []string) {
	a.args = args
	a.config = config
	a.fmt = log.New(os.Stdin, "", 0)
	a.efmt = log.New(os.Stderr, "", 0)
	a.debug = a.getDebug()
	if lang, err := lib.LookupLanguage(config.Get("Language").(string)); err != nil {
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
